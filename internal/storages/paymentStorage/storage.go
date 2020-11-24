package paymentStorage

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Solar-2020/Payment-Backend/pkg/models"
	"strconv"
	"strings"
)

const (
	queryReturningID = "RETURNING id;"
)

var paymentMethodsSchema = map[models.PaymentType]struct{
	tableName string
	cols []string
}{
	models.YoomoneyType: {"yoomoney_account", []string{"account_number"}},
	models.PhoneType: {"phone_payment", []string{"phone_number"}},
	models.CardType: {"bank_card", []string{"bank_title", "bank_logo", "phone_number", "card_number"}},
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *storage {
	return &storage{
		db: db,
	}
}

func (s *storage) InsertPayments(payments []models.Payment, createBy, groupID, postID int) (err error) {
	if len(payments) == 0 {
		return
	}

	const sqlQueryTemplate = `
	INSERT INTO payment(group_id, post_id, create_by, total_cost, payment_account)
	VALUES `

	var params []interface{}

	sqlQuery := sqlQueryTemplate + s.createInsertQuery(len(payments), 5) + queryReturningID

	for i, _ := range payments {
		params = append(params, groupID, postID, createBy, payments[i].TotalCost, payments[i].PaymentAccount)
	}

	for i := 1; i <= len(payments)*5; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return
	}

	rows, err := tx.Query(sqlQuery, params...)
	if err != nil {
		tx.Rollback()
		return err
	}
	i := int(0)
	for rows.Next() {
		rows.Scan(&payments[i].ID)
		i += 1
	}

	for _, payment := range payments {
		err = s.insertMethod(tx, payment.Methods, payment.ID, createBy)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return
}

func (s *storage) SelectPaymentsByPostsIDs(postIDs []int) (payments []models.Payment, err error) {
	payments = make([]models.Payment, 0)
	if len(postIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT p.id, p.total_cost, p.payment_account, p.create_by, p.group_id, p.post_id
	FROM payment AS p
	WHERE p.post_id IN `

	sqlQuery := sqlQueryTemplate + createIN(len(postIDs))

	var params []interface{}

	for i, _ := range postIDs {
		params = append(params, postIDs[i])
	}

	for i := 1; i <= len(postIDs)*1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempPayment models.Payment
		err = rows.Scan(&tempPayment.ID, &tempPayment.TotalCost, &tempPayment.PaymentAccount, &tempPayment.CreateBy, &tempPayment.GroupID, &tempPayment.PostID)
		if err != nil {
			return
		}
		methods, err := s.selectMethods(tempPayment.ID)
		if err != nil {
			return payments, err
		}
		tempPayment.Methods = methods
		payments = append(payments, tempPayment)
	}

	return
}

func (s *storage) SelectPaymentsByPostID(postID int) (payments []models.Payment, err error) {
	postsIDs := make([]int, 1)
	postsIDs = append(postsIDs, postID)
	return s.SelectPaymentsByPostsIDs(postsIDs)
}

func (s *storage) SelectPayment(paymentID int) (payment models.Payment, err error) {
	const sqlQuery = `
	SELECT p.id, p.total_cost, p.payment_account, p.group_id, p.post_id
	FROM payment AS p
	WHERE p.id = $1;`

	err = s.db.QueryRow(sqlQuery, paymentID).Scan(&payment.ID, &payment.TotalCost, &payment.PaymentAccount, &payment.GroupID, &payment.PostID)

	if err != nil {
		return
	}
	methods, err := s.selectMethods(paymentID)
	if err != nil {
		return
	}
	payment.Methods = methods

	return
}


func (s *storage) insertMethod(tx *sql.Tx, methods []models.PaymentMethod, paymentID int, owner int) (err error) {
	if len(methods) == 0 {
		return
	}

	paramsList := map[models.PaymentType][]interface{}{
		models.CardType: make([]interface{}, 0),
		models.PhoneType: make([]interface{}, 0),
		models.YoomoneyType: make([]interface{}, 0),
	}

	const sqlQueryTemplate = `
	INSERT INTO %s(payment_id, owner, %s)
	VALUES `

	for _, method := range methods {
		params := paramsList[method.Type]
		params = append(params, paymentID, owner)
		switch method.Type{
		case models.YoomoneyType:
			params = append(params, method.AccountNumber)
		case models.CardType:
			params = append(params, method.BankName, method.BankLogo, method.PhoneNumber, method.CardNumber)
		case models.PhoneType:
			params = append(params, method.PhoneNumber)
		default:
			return errors.New("bad payment type: " + string(method.Type))
		}
		paramsList[method.Type] = params
	}

	for payType, params := range paramsList {
		if len(params) == 0 {
			continue
		}
		paramsCount := 2 + len(paymentMethodsSchema[payType].cols)

		sqlQuery := fmt.Sprintf(sqlQueryTemplate,
			paymentMethodsSchema[payType].tableName,
			strings.Join(paymentMethodsSchema[payType].cols, ", "),
		)

		sqlQuery = CreatePacketQuery(sqlQuery, paramsCount, len(params)/paramsCount, 0)
		_, err = tx.Exec(sqlQuery, params...)
		if err != nil {
			return err
		}
	}
	return

}

func (s *storage) selectMethods(paymentID int) (methods []models.PaymentMethod, err error) {
	const sqlQueryTemplate = `SELECT owner, payment_id, %s FROM %s WHERE payment_id = $1`

	for payType, schema := range paymentMethodsSchema {
		sqlQuery := fmt.Sprintf(sqlQueryTemplate,
			strings.Join(schema.cols, ", "),
			schema.tableName,
		)

		rows, err := s.db.Query(sqlQuery, paymentID)
		if err != nil {
			return methods, err
		}
		for rows.Next(){
			var tempMethod models.PaymentMethod
			switch payType {
			case models.YoomoneyType:
				err = rows.Scan(&tempMethod.Owner,&tempMethod.PaymentID, &tempMethod.AccountNumber)
				tempMethod.Type = models.YoomoneyType
			case models.CardType:
				err = rows.Scan(&tempMethod.Owner, &tempMethod.PaymentID, &tempMethod.BankName,
					&tempMethod.BankLogo, &tempMethod.PhoneNumber, &tempMethod.CardNumber)
				tempMethod.Type = models.CardType
			case models.PhoneType:
				err = rows.Scan(&tempMethod.Owner,&tempMethod.PaymentID, &tempMethod.PhoneNumber)
				tempMethod.Type = models.PhoneType
			default:
				err = errors.New("bad payment type: " + string(payType))
				return methods, err
			}
			if err != nil {
				return methods, err
			}
			methods = append(methods, tempMethod)
		}
	}
	return
}

func createIN(count int) (queryIN string) {
	queryIN = "("
	for i := 0; i < count; i++ {
		queryIN += "?, "
	}
	queryINRune := []rune(queryIN)
	queryIN = string(queryINRune[:len(queryINRune)-2])
	queryIN += ")"
	return
}

func (s *storage) createInsertQuery(sliceLen int, structLen int) (query string) {
	query = ""
	for i := 0; i < sliceLen; i++ {
		query += "("
		for j := 0; j < structLen; j++ {
			query += "?,"
		}
		// delete last comma
		query = strings.TrimRight(query, ",")
		query += "),"
	}
	// delete last comma
	query = strings.TrimRight(query, ",")

	return
}


// https://github.com/go-park-mail-ru/2019_2_Next_Level/blob/master/pkg/sqlTools/genQuery.go
func CreatePacketQuery(prefix string, batchSize int, batchCount int, initIdx int, postfix ...string) string {
	pack := make([]string, 0, batchCount)
	batch := make([]string, 0, batchSize)

	for i := 0; i < batchCount; i++ {
		for j := 1; j <= batchSize; j++ {
			batch = append(batch, "$"+strconv.Itoa(batchSize*i + j + initIdx))
		}
		pack = append(pack, "(" + strings.Join(batch, ", ") + ")" )
		batch = batch[:0]
	}

	var res bytes.Buffer
	res.WriteString(prefix)
	if prefix[len(prefix)-1] != ' ' {
		res.WriteString(" ")
	}
	res.WriteString(strings.Join(pack, ", "))
	if len(postfix) > 0 {
		res.WriteString(" ")
		res.WriteString(strings.Join(postfix, " "))
	}
	res.WriteString(";")
	return res.String()
}