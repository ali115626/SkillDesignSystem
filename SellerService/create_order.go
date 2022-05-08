package SellerService

import (
	"SeckillDesign/constant"
	//"SeckillDesign/consumer"
	"SeckillDesign/mqConsumer"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
)

//func CreateOrder(){
////	todo 消费MQ
//	message :=consumer.ReceiveMessage()
////	TODO 这个你要unMarshal一下吗？
//	fmt.Println(message)
////	TODO 把订单信息insert到数据库
//}

//func BuildOrder(w http.ResponseWriter, r *http.Request){
//	//	把activityId 和 userID 给传进来
//	//	把order 写入Mq 中
//	err := r.ParseForm()
//	if err != nil {
//		return
//	}
//	requestMap := r.Form
//	activityId := requestMap["activityId"][0]
//	//TODO 这里面加一个用户是不是之前已经购买过的代码逻辑
//	//fmt.Println("userID=",commidtyId)
//	//userId其实是这样给传过来的 不是从请求中获取得到的
//	userId := requestMap["userId"][0]

//	todo 消费MQ

//todo  从messageQUeue中读取订单的信息

//type Order struct{
//	OrderId string
//	UserId string
//	ActivityName string
//	OrderPrice string
//	ActivityId string
//}
//type OrderInfo struct {
//	OrderId      string
//	UserId       string
//	ActivityName string
//	OrderPrice   string
//	ActivityId   string
//	Status       int
//}

func BuildOrderConsumer() {

	//这个要一直for 循环

	//input : activityId
	//        userId
	//
	queueName :="orderMessage"
	message := mqConsumer.ReceiveMessageNormalConsumer(queueName)

	//	TODO 这个你要unMarshal一下吗？
	//fmt.Println(message)
	//todo  检查一下  库存里面有没有
	//TODO 有的话  给前端返回1 ：创建订单成功  2：已经支付了（支付成功了）

	//	TODO 把订单信息insert到数据库
	orderInfo := &constant.OrderInfo{}

	err := json.Unmarshal([]byte(message), orderInfo)
	if err != nil {
		fmt.Println("json unmarshal orderInfo err,err=", err)
	}
	//fmt.Println(orderInfo)
	//fmt.Printf("orderInfo=%v\n",orderInfo)

	db, err := sql.Open("mysql", "root:123456@/seckill_scheme?charset=utf8")
	if err != nil {
		fmt.Println("open database error,err=", err)
	}

	//var stocks string
	//todo  数据库这里 你就不能用一下乐观锁吗？  update ActivityTable set available_stock=available_stock-1,locked_stock=locked_stock+1 WHERE activityId= orderinfo.activityId and stocks

	//todo 这个万一fail了 库存扣减失败  status = 0  这里的话 就直接return
	//todo 这个如果成功了的话   库存扣减成功 status=1

	//TODO	先去查库存  看看还有没有这个订单信息  如果有库存  把status改成 1  否则改成 status = 0
	//	err = db.QueryRow("SELECT stocks FROM ActivityTable WHERE activityId=?", orderInfo.ActivityId).Scan(&stocks)
	//	//发现即使db.QueryRow(）这里面ELECT delete_status FROM blog_info WHERE title_id SQL语句出问题了  也不会报错
	//	if err != nil {
	//	_, err = db.Exec("update ActivityTable set available_stock=available_stock-1,locked_stock=locked_stock+1 WHERE activityId= and available_stock > ) values(?,?)",orderInfo.ActivityId,0)
	//	if err != nil {
	//		fmt.Println("exec failed, err=", err)
	//		//	todo 我就奇怪了   万一我在这边一直点下去怎么办   代码层校验是否重复上传   怎么办呢
	//	}

	//result, err := db.Exec("update ActivityTable set available_stock=available_stock-1,locked_stock=locked_stock+1 WHERE activityId=? and available_stock > ?) values(?,?)",ActivityId,0)
	//orderInfo.ActivityId = "18"
	result, err := db.Exec("update ActivityTable set available_stock=available_stock-1,locked_stock=locked_stock+1 WHERE activityId=? and available_stock > ?", orderInfo.ActivityId, 0)

	if err != nil {
		fmt.Println("exec failed, err=", err)
		//	todo 我就奇怪了   万一我在这边一直点下去怎么办   代码层校验是否重复上传   怎么办呢
	}

	RowsAffected, err := result.RowsAffected()

	if RowsAffected == 0 {
		//Status
		orderInfo.Status = 0
		//	TODO status :将其设为  库存扣减失败
		return //前端发现你这个东西等于0了   然后就开始  todo 这样的话  也没有必要扣减数据库了
	}
	//
	orderInfo.Status = 1

	//TODO 这个要做一下转化哈

	userId, err := strconv.Atoi(orderInfo.UserId)
	if err != nil {
		fmt.Println(err)
	}

	orderPrice, err := strconv.Atoi(orderInfo.OrderPrice)
	if err != nil {
		fmt.Println(err)

	}

	activityId, err := strconv.Atoi(orderInfo.ActivityId)
	if err != nil {
		fmt.Println(err)

	}

	status := strconv.Itoa(orderInfo.Status)

	inserResult, err := db.Exec("insert into OrderInfoTable(orderId,userId,activityName,orderPrice,activityId,status) values(?,?,?,?,?,?)", orderInfo.OrderId, userId, orderInfo.ActivityName, orderPrice, activityId, status)
	if err != nil {
		fmt.Println("exec failed, err=", err)
		//	todo 我就奇怪了   万一我在这边一直点下去怎么办   代码层校验是否重复上传   怎么办呢
	}

	//todo  这样的话  insert进数据库就成功了

	fmt.Println(inserResult.RowsAffected())

	//TODO  否则 insert到数据库中呗

	//	TODO  如果扣减失败了  直接return回去  扣减成功  修改status insert进数据库

	//		fmt.Println(err)
	//		fmt.Println("select  paper_content  error")
	//		return
	//
	//	}
	//
	//	if stocks > 0{
	//
	//	}

	//todo  所以 你那边还要去确认一下   有没有库存

}
