package paymentStorage

import (
	"database/sql"
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

func (s *storage) InsertPayments(payments []Payment, createBy, groupID, postID int) (err error) {
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

func (s *storage) SelectPaymentsByPostsIDs(postIDs []int) (payments []Payment, err error) {
	payments = make([]Payment, 0)
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
		var tempPayment Payment
		err = rows.Scan(&tempPayment.ID, &tempPayment.TotalCost, &tempPayment.PaymentAccount, &tempPayment.CreateBy, &tempPayment.GroupID, &tempPayment.PostID)
		if err != nil {
			return
		}
		payments = append(payments, tempPayment)
	}

	return
}

func (s *storage) SelectPaymentsByPostID(postID int) (payments []Payment, err error) {
	postsIDs := make([]int, 1)
	postsIDs = append(postsIDs, postID)
	return s.SelectPaymentsByPostsIDs(postsIDs)
}

func (s *storage) SelectPayment(paymentID int) (payment Payment, err error) {
	const sqlQuery = `
	SELECT p.id, p.total_cost, p.payment_account, p.group_id, p.post_id
	FROM payment AS p
	WHERE p.id = $1;`

	err = s.db.QueryRow(sqlQuery, paymentID).Scan(&payment.ID, &payment.TotalCost, &payment.PaymentAccount, &payment.GroupID, &payment.PostID)

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
