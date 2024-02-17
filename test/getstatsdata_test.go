package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/psyark/estat"
)

func init() {
	if err := godotenv.Load("testdata/secret.env"); err != nil {
		panic(err)
	}
}

func TestGetStatsData(t *testing.T) {
	t.Skip()
	ctx := context.Background()

	{
		// 件数取得
		query := url.Values{}
		query.Set("appId", os.Getenv("appId"))
		query.Set("limit", "1")

		content, err := estat.GetStatsList(ctx, query)
		if err != nil {
			t.Fatal(err)
		}

		// ランダムな位置で10000件取得
		const windowSize = 10000
		query.Set("startPosition", strconv.Itoa(rand.Intn(content.DatalistInf.Number-windowSize)))
		query.Set("limit", strconv.Itoa(windowSize))

		content, err = estat.GetStatsList(ctx, query)
		if err != nil {
			t.Fatal(err, query) // 偶発的なエラーのためにqueryを見る
		}

		// そこからランダムに取得
		for i := 0; i < 10; i++ {
			statsDataId := content.DatalistInf.TableInf[rand.Intn(len(content.DatalistInf.TableInf))].ID
			query = url.Values{}
			query.Set("appId", os.Getenv("appId"))
			query.Set("statsDataId", statsDataId)
			_, err = estat.GetStatsData(ctx, query, estat.WithDataHandler(func(data []byte) error {
				return os.WriteFile(fmt.Sprintf("testdata/%s.json", statsDataId), data, 0666)
			}))
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") && entry.Name() != "list.json" {
			statsDataId := strings.TrimSuffix(entry.Name(), ".json")
			t.Run(statsDataId, func(t *testing.T) {
				query := url.Values{}
				query.Set("appId", os.Getenv("appId"))
				query.Set("statsDataId", statsDataId)

				jsonFile := fmt.Sprintf("testdata/%s.json", statsDataId)

				if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
					_, err := estat.GetStatsData(ctx, query, estat.WithDataHandler(func(data []byte) error {
						return os.WriteFile(fmt.Sprintf("testdata/%s.json", statsDataId), data, 0666)
					}))
					if err != nil {
						t.Fatal(err)
					}
				}

				data, err := os.ReadFile(jsonFile)
				if err != nil {
					t.Fatal(err)
				}

				err = cyclicTest(data, func() ([]byte, error) {
					typed := estat.GetStatsDataContainer{}
					if err := json.Unmarshal(data, &typed); err != nil {
						return nil, err
					}
					return json.Marshal(typed)
				})
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func Example_s0003281511() {
	result := map[string]string{}
	s("0003281511", func(gsd *estat.GetStatsDataContent, v estat.Value) {
		if v.Time == "2018000000" {
			x := gsd.StatisticalData.ClassInf.Cat04().GetClass(v.Cat04)
			result[x.Name] = fmt.Sprintf("%v%v", v.Value, v.Unit)
		}
	})

	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))

	// output:
	// 警察庁 交通事故の発生状況
	// 交通事故の発生状況 原付以上運転者（第1当事者）の年齢層別免許保有者10万人当たり交通事故件数の推移
	// tab 表章項目
	//     3430 = 免許保有者10万人当たり事故件数
	// cat04 年齢層
	//     110 = 全年齢層
	//     150 = 15歳以下
	//     170 = 16～19歳
	//     180 = 20～24歳
	//     190 = 25～29歳
	//     200 = 30～34歳
	//     210 = 35～39歳
	//     220 = 40～44歳
	//     230 = 45～49歳
	//     240 = 50～54歳
	//     250 = 55～59歳
	//     260 = 60～64歳
	//     270 = 65～69歳
	//     280 = 70～74歳
	//     290 = 75～79歳
	//     300 = 80～84歳
	//     310 = 85歳以上
	//     320 = (再掲)16～24歳
	//     340 = (再掲)65歳以上
	//     350 = (再掲)70歳以上
	//     360 = (再掲)75歳以上
	//     370 = (再掲)80歳以上
	// time 時間軸（年次）
	//     2018000000 = 2018年
	//     2017000000 = 2017年
	//     2016000000 = 2016年
	//     2015000000 = 2015年
	//     2014000000 = 2014年
	//     2013000000 = 2013年
	//     2012000000 = 2012年
	//     2011000000 = 2011年
	//     2010000000 = 2010年
	//     2009000000 = 2009年
	//     2008000000 = 2008年
	//     2007000000 = 2007年
	//     2006000000 = 2006年
	//     2005000000 = 2005年
	// {
	//   "(再掲)16～24歳": "973.1件",
	//   "(再掲)65歳以上": "483.3件",
	//   "(再掲)70歳以上": "512.4件",
	//   "(再掲)75歳以上": "566.4件",
	//   "(再掲)80歳以上": "615.7件",
	//   "15歳以下": "-件",
	//   "16～19歳": "1489.2件",
	//   "20～24歳": "876.9件",
	//   "25～29歳": "624.0件",
	//   "30～34歳": "487.5件",
	//   "35～39歳": "433.7件",
	//   "40～44歳": "432.2件",
	//   "45～49歳": "431.7件",
	//   "50～54歳": "414.0件",
	//   "55～59歳": "415.6件",
	//   "60～64歳": "426.4件",
	//   "65～69歳": "438.4件",
	//   "70～74歳": "458.6件",
	//   "75～79歳": "533.3件",
	//   "80～84歳": "604.5件",
	//   "85歳以上": "645.9件",
	//   "全年齢層": "494.1件"
	// }
}

func Example_s0004009602() {
	result := map[string]string{}
	s("0004009602", func(gsd *estat.GetStatsDataContent, v estat.Value) {
		if v.Tab == "0060" && v.Cat01 == "6" {
			time := v.TimeClass(gsd.StatisticalData.ClassInf)
			result[time.Name] = fmt.Sprintf("%v%v", v.Value, v.Unit)
		}
	})

	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))

	// output:
	// 国税庁 民間給与実態統計 結果表（新たな復元推計手法により計算）
	// map[$:全国計表　第1表　給与所得者数・給与額・税額　業種別　（2014年～） @no:00103]
	// tab 表章項目
	//     0010 = 給与所得者数(３月末)
	//     0020 = 給与所得者数(６月末)
	//     0030 = 給与所得者数(９月末)
	//     0040 = 給与所得者数(１２月末)
	//     0050 = 給与所得者数(年間月平均)
	//     0060 = 給与額(総額)
	//     0070 = 給与額(平均)
	//     0080 = 税額(総額)
	//     0090 = 税額(平均)
	// cat01 業種（2008年～）
	//     1 = 建設業
	//     2 = 製造業
	//     3 = 卸売業，小売業
	//     4 = 宿泊業，飲食サービス業
	//     5 = 金融業，保険業
	//     6 = 不動産業，物品賃貸業
	//     7 = 運輸業，郵便業
	//     8 = 電気･ガス･熱供給・水道業
	//     9 = 情報通信業
	//     10 = 学術研究，専門・技術サービス業、教育，学習支援業
	//     11 = 医療，福祉
	//     12 = 複合サービス事業
	//     13 = サービス業
	//     14 = 農林水産・鉱業
	//     15 = 合計
	// time 年
	//     2022000000 = 2022年
	//     2021000000 = 2021年
	//     2020000000 = 2020年
	//     2019000000 = 2019年
	//     2018000000 = 2018年
	//     2017000000 = 2017年
	//     2016000000 = 2016年
	//     2015000000 = 2015年
	//     2014000000 = 2014年
	// {
	//   "2014年": "4699923百万円",
	//   "2015年": "4449083百万円",
	//   "2016年": "4872191百万円",
	//   "2017年": "5070325百万円",
	//   "2018年": "5342272百万円",
	//   "2019年": "5131226百万円",
	//   "2020年": "5080431百万円",
	//   "2021年": "5114356百万円",
	//   "2022年": "5787558百万円"
	// }
}

func Example_s0003354197() {
	result := map[string]map[string]string{}
	s("0003354197", func(gsd *estat.GetStatsDataContent, v estat.Value) {
		if v.Cat03 == "100" && v.Time == "2001100000" {
			cat01 := v.Cat01Class(gsd.StatisticalData.ClassInf)
			cat02 := v.Cat02Class(gsd.StatisticalData.ClassInf)
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

func Example_s0002019042() {
	result := map[string]string{}
	s("0002019042", func(gsd *estat.GetStatsDataContent, v estat.Value) {
		if v.Cat01 == "1001" {
			cat02 := v.Cat02Class(gsd.StatisticalData.ClassInf)
			result[cat02.Name] = fmt.Sprintf("%v%v", v.Value, v.Unit)
		}
	})

	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))

	// output:
	// 農林水産省 新規就農者調査 確報 令和３年新規就農者調査結果
	// map[$:新規就農者調査 就農形態別新規就農者数 @no:1]
	// cat01 (B002-03-2-001)男女年齢別
	//     1001 = 男女計
	//     1002 = 男女計_49歳以下
	//     1003 = 男女計_49歳以下_15歳～19歳
	//     1004 = 男女計_49歳以下_20歳～29歳
	//     1005 = 男女計_49歳以下_30歳～39歳
	//     1006 = 男女計_49歳以下_40歳～49歳
	//     1007 = 男女計_50歳～59歳
	//     1008 = 男女計_60歳～64歳
	//     1009 = 男女計_65歳以上
	//     1010 = 男計
	//     1011 = 男_49歳以下
	//     1012 = 男_49歳以下_15歳～19歳
	//     1013 = 男_49歳以下_20歳～29歳
	//     1014 = 男_49歳以下_30歳～39歳
	//     1015 = 男_49歳以下_40歳～49歳
	//     1016 = 男_50歳～59歳
	//     1017 = 男_60歳～64歳
	//     1018 = 男_65歳以上
	//     1019 = 女計
	//     1020 = 女_49歳以下
	//     1021 = 女_49歳以下_15歳～19歳
	//     1022 = 女_49歳以下_20歳～29歳
	//     1023 = 女_49歳以下_30歳～39歳
	//     1024 = 女_49歳以下_40歳～49歳
	//     1025 = 女_50歳～59歳
	//     1026 = 女_60歳～64歳
	//     1027 = 女_65歳以上
	// cat02 (B002-03-1-001)就農形態別
	//     1001 = 計_平成30年
	//     1002 = 計_令和元年
	//     1003 = 計_令和2年
	//     1004 = 計_令和3年
	//     1005 = 新規自営農業就農者_平成30年
	//     1006 = 新規自営農業就農者_令和元年
	//     1007 = 新規自営農業就農者_令和2年
	//     1008 = 新規自営農業就農者_令和3年
	//     1009 = 新規雇用就農者_平成30年
	//     1010 = 新規雇用就農者_令和元年
	//     1011 = 新規雇用就農者_令和2年
	//     1012 = 新規雇用就農者_令和3年
	//     1013 = 新規参入者_平成30年
	//     1014 = 新規参入者_令和元年
	//     1015 = 新規参入者_令和2年
	//     1016 = 新規参入者_令和3年
	// {
	//   "新規参入者_令和2年": "3580人",
	//   "新規参入者_令和3年": "3830人",
	//   "新規参入者_令和元年": "3200人",
	//   "新規参入者_平成30年": "3240人",
	//   "新規自営農業就農者_令和2年": "40100人",
	//   "新規自営農業就農者_令和3年": "36890人",
	//   "新規自営農業就農者_令和元年": "42740人",
	//   "新規自営農業就農者_平成30年": "42750人",
	//   "新規雇用就農者_令和2年": "10050人",
	//   "新規雇用就農者_令和3年": "11570人",
	//   "新規雇用就農者_令和元年": "9940人",
	//   "新規雇用就農者_平成30年": "9820人",
	//   "計_令和2年": "53740人",
	//   "計_令和3年": "52290人",
	//   "計_令和元年": "55870人",
	//   "計_平成30年": "55810人"
	// }
}

func s(statsDataId string, f func(*estat.GetStatsDataContent, estat.Value)) {
	ctx := context.Background()

	query := url.Values{}
	query.Set("appId", os.Getenv("appId"))
	query.Set("statsDataId", statsDataId)

	gsd, err := estat.GetStatsData(ctx, query)
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
