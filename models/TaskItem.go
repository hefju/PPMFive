package models

import ("time"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type TaskItem struct {
	Id int64            `json:"id"` //主键
	UserId int64         //用户ID
	DateOfTask time.Time //任务的日期
	Title string          `json:"title"`//任务的标题
//	Desc string          //任务的描述
	Done bool          `json:"done"`  //是否完成
//	Created time.Time  `xorm:CREATED`   //添加时间
//	Updated time.Time    //最后更新时间
	Version int     `xorm:"version"`     //版本
}

// ORM 引擎
var x *xorm.Engine
func init() {
	// 创建 ORM 引擎与数据库
	var err error
	x, err = xorm.NewEngine("sqlite3", "./test.db")
	if err != nil {
		log.Fatalf("Fail to create engine: %v\n", err)
	}
	x.SetTableMapper(core.SameMapper{})
	x.SetMapper(core.SameMapper{})

	x.ShowSQL = true
	x.ShowInfo = true
	x.ShowDebug = true
	x.ShowErr = true
	x.ShowWarn = true
	// 同步结构体与数据表
	if err = x.Sync(new(TaskItem)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}
}

func AddTask(t *TaskItem) (int64,error){
	affected,err:=x.Insert(t)
	return affected,err
}

func UpdateTask(t *TaskItem) (int64,error){
	affected,err:=x.Id(t.Id).Cols("done","title").Update(t)//为什么只能用cols才能更新成功
	return affected,err
}
func DeleteTask(id int64) (int64,error){
	affected,err:=x.Id(id).Delete(new(TaskItem))
	return affected,err
}

func GetTaskByID(id int64) (*TaskItem,error){
	task:=new(TaskItem)
	_,err:=x.Id(id).Get(task)
	return task,err
}

func GetTaskList(date time.Time)([]*TaskItem,error){
	tasks := make([]*TaskItem, 0)
	err := x.Find(&tasks)//.Cols("id","title","done") //x.Where("date_of_task=?",date).Find(&tasks)
	if err!=nil{
		log.Println("GetTaskList: ",err)
	}
	return tasks, err
}




