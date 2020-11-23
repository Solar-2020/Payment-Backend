package paymentStorage

import (
	"database/sql"
	models2 "github.com/Solar-2020/Payment-Backend/internal/models"
	"github.com/Solar-2020/Payment-Backend/pkg/models"
	"strconv"
	"strings"
)

const (
	queryReturningID = "RETURNING id;"
)

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

	_, err = s.db.Exec(sqlQuery, params...)

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

	return
}

func (s *storage) SelectPaids(paymentID int) (paids []models2.Paid, err error) {
	const sqlQuery = `
	SELECT p.id, p.user_id, p.cost, p.paid_at, p.requisite_type_id, p.requisite_id
	FROM paid AS p
	WHERE p.payment_id = $1;`

	rows, err := s.db.Query(sqlQuery, paymentID)
	if err != nil {
		return
	}

	for rows.Next() {
		var tempPaid models2.Paid
		err = rows.Scan(&tempPaid.PaidID, &tempPaid.PayerID, &tempPaid.Cost, &tempPaid.PaidAt, &tempPaid.RequisiteType, &tempPaid.RequisiteID)
		if err != nil {
			return
		}
		paids = append(paids, tempPaid)
	}

	for i, _ := range paids {
		switch paids[i].RequisiteType {
		case 1:
			bankCard, err := s.selectBankRequisite(paids[i].RequisiteID)
			if err != nil {
				return
			}
			paids[i].BankCard = &bankCard
		case 2:
			phonePayment, err := s.selectPhoneRequisite(paids[i].RequisiteID)
			if err != nil {
				return
			}
			paids[i].PhonePayment = &phonePayment
		case 3:
			youMoneyAccount, err := s.selectYouMoneyRequisite(paids[i].RequisiteID)
			if err != nil {
				return
			}
			paids[i].YouMoneyAccount = &youMoneyAccount
		}
	}
	return
}

func (s *storage) selectBankRequisite(bankRequisiteID int) (bankCard models2.BankCard, err error) {
	const sqlQuery = `
	SELECT bc.id, bc.bank_title, bc.phone_number, bc.card_number, bc.owner
	FROM bank_card AS bc
	WHERE bc.id = $1;`

	err = s.db.QueryRow(sqlQuery, bankRequisiteID).Scan(&bankCard.ID, &bankCard.BankTitle, &bankCard.PhoneNumber, &bankCard.CardNumber, &bankCard.Owner)
	return
}

func (s *storage) selectPhoneRequisite(bankRequisiteID int) (phonePayment models2.PhonePayment, err error) {
	const sqlQuery = `
	SELECT pp.id, pp.phone_number, pp.owner
	FROM phone_payment AS pp
	WHERE pp.id = $1;`

	err = s.db.QueryRow(sqlQuery, bankRequisiteID).Scan(&phonePayment.ID, &phonePayment.PhoneNumber, &phonePayment.Owner)
	return
}

func (s *storage) selectYouMoneyRequisite(bankRequisiteID int) (youMoneyAccount models2.YouMoneyAccount, err error) {
	const sqlQuery = `
	SELECT ya.id, ya.account_number, ya.owner
	FROM yoomoney_account AS ya
	WHERE ya.id = $1;`

	err = s.db.QueryRow(sqlQuery, bankRequisiteID).Scan(&youMoneyAccount.ID, &youMoneyAccount.AccountNumber, &youMoneyAccount.Owner)
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
