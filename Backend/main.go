package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)



func StoreWord(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}


	// 解析JSON
	WordData := make(map[string]interface{})
	err = json.Unmarshal(body, &WordData)
	if err != nil {
		fmt.Println(err)
		return
	}
	word := WordData["word"]
	fmt.Println("接收到的单词是：",word)

	s := `select from "LE_word" where word = $1`
	err = Connection.QueryRow(context.Background(), s, word).Scan()
	// 单词不存在
	if err == pgx.ErrNoRows {
		s = `insert into "LE_word"(word, create_time) values($1, $2)`
		_, err = Connection.Exec(context.Background(), s, word, time.Now())
		if err != nil {
			fmt.Println("insert word err:", err)
		}
	} else {
		s = `update "LE_word" set repetition=repetition+1 where word=$1`
		_, err = Connection.Exec(context.Background(), s, word)
		if err != nil {
			fmt.Println("update word failed:", err)
		}
	}
}
func AddSource(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("read body err:", err)
		return
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("json Unmarshal err:", err)
		return
	}

	title := data["title"]
	article := data["article"]
	fmt.Println(title)
	fmt.Println(article)
	s := `insert into "LE_article"(title, article) values($1, $2)`
	_, err = Connection.Exec(context.Background(), s, title, article)
	if err != nil {
		fmt.Println("insert article err:", err)
		return
	}
}
func GetArticleList(w http.ResponseWriter, r *http.Request) {
	s := `select title from "LE_article"`
	rows, err := Connection.Query(context.Background(), s)
	if err != nil {
		fmt.Println("select title err:", err)
		return
	}
	var titles []string
	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		if err != nil {
			fmt.Println("Scan err:", err)
			return
		}
		titles = append(titles, title)
	}
	fmt.Println("titles", titles)
	titlesJson, err := json.Marshal(titles)
	if err != nil {
		fmt.Println("Marshal err: ", err)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"data":%s}`, titlesJson)))
}
func GetArticle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	fmt.Println("title", title)
	s := `select article from "LE_article" where title = $1`
	var article string
	err := Connection.QueryRow(context.Background(), s, title).Scan(&article)
	if err != nil {
		fmt.Println("select err:", err)
		return
	}
	fmt.Println(article)
	articleJson, err := json.Marshal(article)
	if err != nil {
		fmt.Println("Marshal err: ", err)
	}
	w.Write(articleJson)
}
func GetLyric(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("./lyric.lrc")
	if err != nil {
		fmt.Println(err)
		return
	}
	lyricJson, err := json.Marshal(string(file))
	if err != nil {
		fmt.Println("Marshal err: ", err)
	}
	w.Write(lyricJson)
}
var Connection *pgxpool.Pool
func main() {
	config, err := pgxpool.ParseConfig("postgres://postgres:123456789@localhost:1433/postgres")
	//config, err := pgxpool.ParseConfig("postgres://postgres:123456789@LearnE:5432/postgres")
	if err != nil {
		fmt.Println("err1:",err)
	}
	Connection, err = pgxpool.ConnectConfig(context.Background(), config)
	//处理连接失败的情况
	if err != nil {
		fmt.Println("err2:",err)
	}
	// 定义 HTTP 路由
	http.HandleFunc("/api/reading/storeWord",StoreWord)
	http.HandleFunc("/api/reading/addsource",AddSource)
	http.HandleFunc("/api/reading/getArticleList",GetArticleList)
	http.HandleFunc("/api/reading/getArticle", GetArticle)
	http.HandleFunc("/api/reading/geiLyric",GetLyric)
	// 启动 HTTP 服务器
	fmt.Println("Starting server on :1234")
	log.Fatal(http.ListenAndServe(":1235", nil))
}
