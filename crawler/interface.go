package crawler

type ICrawler[T any] interface {
	// SetAddress 設定爬蟲的網址
	SetAddress(address string)

	// HandleRawDataToExpected 處理原始資料轉換為預期的資料
	HandleRawDataToExpected() (expected T, err error)

	// SaveToStorage 儲存資料到指定的儲存空間
	SaveToStorage(expected T) error
}
