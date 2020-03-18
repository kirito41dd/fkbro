package data

type Alert struct {
	ID 				int64
	TimeStamp		int64
	Symbol 			string
	Hash 			string
	Amount 			int64
	AmountUsd		int64
	FromAddr 		string
	FromOwner		string
	TomAddr 		string
	TomOwner		string
}

func (al *Alert) Insert() error {
	statement := "insert into alert (time_stamp,symbol,hash,amount,amount_usd,from_addr,from_owner,to_addr,to_owner) values(?,?,?,?,?,?,?,?,?)"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		Log.Debug(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(al.TimeStamp,al.Symbol,al.Hash,al.Amount,al.AmountUsd,al.FromAddr,al.FromOwner,al.TomAddr, al.TomOwner)
	if err != nil {
		if err != nil {
			Log.Debug(err)
			return err
		}
	}
	return nil
}

func (al *Alert) CalcWarnLevel() string {
	n := int(al.AmountUsd / 10000000)
	var s string
	if n > 10 {
		n = 10
	}
	for i := 0; i < n; i++ {
		s += "🚨"
	}
	return s
}

func GetAlertsByTimeStamp(symbol string, timestamp int64) []*Alert {
	statement := "select id,time_stamp,symbol,hash,amount,amount_usd,from_addr,from_owner,to_addr,to_owner from alert where symbol=? and time_stamp>=? order by time_stamp"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		Log.Debug(err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(symbol, timestamp)
	if err != nil {
		Log.Debug(err)
		return nil
	}
	defer rows.Close()
	re := make([]*Alert,0)
	for rows.Next() {
		al := &Alert{}
		err = rows.Scan(&al.ID, &al.TimeStamp, &al.Symbol, &al.Hash, &al.Amount, &al.AmountUsd, &al.FromAddr, &al.FromOwner, &al.TomAddr, &al.TomOwner)
		if err != nil {
			Log.Debug(err)
			return nil
		}
		re = append(re, al)
	}
	return re
}

func DeleteAlert(timestamp int64) {
	statement := "delete from alert where time_stamp < ?"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		Log.Debug(err)
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(timestamp)
	if err != nil {
		Log.Debug(err)
		return
	}
	rows, _ := res.RowsAffected()
	Log.Debug("delete success, rows affect", rows)
}

