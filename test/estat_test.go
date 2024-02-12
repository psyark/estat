package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/psyark/estat"
	"github.com/wI2L/jsondiff"
)

func init() {
	if err := godotenv.Load("testdata/secret.env"); err != nil {
		panic(err)
	}
}

func TestXxx(t *testing.T) {
	for _, statsDataId := range []string{"0003354197", "0004009602", "0003313482", "0002019042"} {
		statsDataId := statsDataId
		t.Run(statsDataId, func(t *testing.T) {
			query := url.Values{}
			query.Set("appId", os.Getenv("appId"))
			query.Set("statsDataId", statsDataId)

			resp, err := http.Get("http://api.e-stat.go.jp/rest/3.0/app/json/getStatsData?" + query.Encode())
			if err != nil {
				t.Fatal(err)
			}

			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			untyped := map[string]any{}
			if err := json.Unmarshal(data, &untyped); err != nil {
				t.Fatal(err)
			}

			typed := estat.Response{}
			if err := json.Unmarshal(data, &typed); err != nil {
				t.Fatal(err)
			}

			data1, _ := json.MarshalIndent(untyped, "", "  ")

			os.WriteFile(fmt.Sprintf("testdata/%s.json", statsDataId), data1, 0666)

			data2, _ := json.MarshalIndent(typed, "", "  ")

			patch, err := jsondiff.CompareJSON(data1, data2)
			if err != nil {
				t.Fatal(err)
			}

			for _, op := range patch {
				fmt.Printf("%v %v: value=%v, oldValue=%v\n", op.Type, op.Path, op.Value, op.OldValue)
			}

			if len(patch) != 0 {
				t.Fatal("unmatch")
			}
		})
	}
}

func Example_s0003354197() {
	result := map[string]map[string]string{}
	s("0003354197", func(gsd *estat.GetStatsData, v estat.Value) {
		if v.Cat03 == "100" && v.Time == "2001100000" {
			cat01 := getClassName(gsd.StatisticalData, "cat01", v.Cat01)
			cat02 := getClassName(gsd.StatisticalData, "cat02", v.Cat02)
			if result[cat02.Name] == nil {
				result[cat02.Name] = map[string]string{}
			}
			result[cat02.Name][fmt.Sprintf("%s %s", cat01.Code, cat01.Name)] = fmt.Sprintf("%v%v", v.Value, v.Unit)
		}
	})

	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))

	// output:
	// こども家庭庁 情報化社会と青少年に関する調査 青少年に対する質問
	// 青少年 Q18　[カード16]　インターネットを通じて、実際に利用したことのあるものはどれですか。次の中からあてはまるものをいくつでもあげてください。(M.A.)　性・年齢1
	// cat01 利用内容
	//     110 = 該当数
	//     120 = 企業・政府・団体のホームページ
	//     130 = 個人のホームページを見る
	//     140 = インターネット・バンキング
	//     150 = オンラインゲーム
	//     160 = インターネット・ショッピング
	//     170 = 航空券、演劇のチケット予約・購入
	//     180 = インターネット・オークション
	//     190 = 動画・音楽などのダウンロード
	//     200 = 株取引などの金融取引
	//     210 = メールマガジン
	//     220 = 出会い系サイト
	//     230 = その他
	//     240 = 無回答
	//     250 = 回答計
	// cat02 年齢
	//     100 = 総数
	//     110 = 12～14歳
	//     120 = 15～17歳
	//     140 = 18～22歳
	//     150 = 23～30歳
	// cat03 性別
	//     100 = 総数
	//     110 = 男性
	//     120 = 女性
	// time 時間軸(年度次)
	//     2001100000 = 2001年度
	// {
	//   "総数": {
	//     "110 該当数": "1996人",
	//     "120 企業・政府・団体のホームページ": "52.4%",
	//     "130 個人のホームページを見る": "53.7%",
	//     "140 インターネット・バンキング": "6.1%",
	//     "150 オンラインゲーム": "20.1%",
	//     "160 インターネット・ショッピング": "25.1%",
	//     "170 航空券、演劇のチケット予約・購入": "17.6%",
	//     "180 インターネット・オークション": "15.2%",
	//     "190 動画・音楽などのダウンロード": "47.1%",
	//     "200 株取引などの金融取引": "1.4%",
	//     "210 メールマガジン": "20.3%",
	//     "220 出会い系サイト": "5.3%",
	//     "230 その他": "0.9%",
	//     "240 無回答": "1.9%",
	//     "250 回答計": "267.1%"
	//   }
	// }
}

func s(statsDataId string, f func(*estat.GetStatsData, estat.Value)) {
	query := url.Values{}
	query.Set("appId", os.Getenv("appId"))
	query.Set("statsDataId", statsDataId)

	gsd, err := estat.GetStats(query)
	if err != nil {
		panic(err)
	}

	tableInf := gsd.StatisticalData.TableInf
	fmt.Println(tableInf.GovOrg.Annotation, tableInf.StatisticsName)
	fmt.Println(tableInf.Title)

	for _, classObj := range gsd.StatisticalData.ClassInf.ClassObj {
		fmt.Println(classObj.ID, classObj.Name)
		for _, class := range classObj.Class {
			fmt.Printf("    %s = %s\n", class.Code, class.Name)
		}
	}

	for _, v := range gsd.StatisticalData.DataInf.Value {
		f(gsd, v)
	}
}

func getClassName(statisticalData estat.StatisticalData, classID string, classCode string) *estat.Class {
	for _, classObj := range statisticalData.ClassInf.ClassObj {
		if classObj.ID == classID {
			for _, class := range classObj.Class {
				if class.Code == classCode {
					return &class
				}
			}
		}
	}
	return nil
}
