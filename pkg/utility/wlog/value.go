package wlog

import (
	"encoding/json"
	"log/slog"
	"unsafe"
)

// JsonValue 將資料轉換為 slog.Value，並根據 asString 決定輸出形式。
//
// 參數：
//  - asString: 將資料表示為 string 形式
//  - data: 只支援 struct, []byte, map, slice, 不支援基本型別
//
// Note:
//
// 除非重寫 Handler, 目前找不到一個簡單的方法, 同時滿足 JsonHandler 以及 TextHandler 都是完美的 json 字串.
//
// asString 參數決定資料輸出的形式，根據不同情境需要謹慎選擇：
//
// 1. when asString = true
//      - config.SetJsonFormat(true)，資料會被包裹在引號中
//      {"time":"2025-01-10T23:15:45+08:00","level":"INFO","msg":"","struct":"{\"name\":\"caesar\"}","no":9487}
//
//      - config.SetJsonFormat(false)，資料會被包裹在引號中
//      2025-01-10T23:15:10+08:00 INF  struct="{\"name\":\"caesar\"}" no=9487
//
// 2. when asString = false
//      - config.SetJsonFormat(true)，資料保持良好的結構化表現
//      {"time":"2025-01-10T23:16:07+08:00","level":"INFO","msg":"","struct":{"name":"caesar"},"no":9487}
//
//      - config.SetJsonFormat(false)，幾乎不可讀，應避免使用
//      2025-01-10T23:16:32+08:00 INF  struct="[123 34 110 97 109 101 34 58 34 99 97 101 115 97 114 34 125]" no=9487
func JsonValue(asString bool, data any) slog.Value {
	var bData []byte
	var err error

	switch v := data.(type) {
	case []byte:
		bData = v
	default:
		bData, err = json.Marshal(data)
		if err != nil {
			return slog.StringValue(err.Error())
		}
	}

	if asString {
		return slog.StringValue(unsafe.String(unsafe.SliceData(bData), len(bData)))
	}
	return slog.AnyValue(json.RawMessage(bData))
}

func PtrValue[T any](val *T) slog.Value {
	if val == nil {
		return slog.StringValue("nil")
	}
	return slog.AnyValue(*val)
}
