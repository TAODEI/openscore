package controllers

import (
	"encoding/json"
	"openscore/models"
	"openscore/requests"
	"openscore/responses"
	"strconv"
	"strings"
)


/**
 9.大题选择列表
 */
func (c *SupervisorApiController) QuestionList() {
	defer c.ServeJSON()
	var requestBody requests.QuestionList
	var resp Response
	var  err error

	err =json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	//supervisorId := requestBody.SupervisorId
	//----------------------------------------------------
	//获取大题列表
	topics  := make([]models.Topic,0)
	err = models.GetTopicList(&topics)
    if err!=nil {
    	resp  = Response{"20000","GetTopicList err ",err}
		c.Data["json"] = resp
		return
	}

	var questions = make([]responses.QuestionListVO,len(topics))
	for i := 0; i < len(topics); i++ {

		questions[i].QuestionId=topics[i].Question_id
		questions[i].QuestionName=topics[i].Question_name

	}

	//----------------------------------------------------
	data := make(map[string]interface{})
	data["questionsList"] =questions
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp
}

/**
 10.用户登入信息表
 */
func (c *SupervisorApiController) UserInfo() {
	defer c.ServeJSON()
	var requestBody requests.UserInfo
	var resp  Response
	var err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	supervisorId := requestBody.SupervisorId

//----------------------------------------------------
    user := models.User{User_id: supervisorId}
	err = user.GetUser(supervisorId)
	if err!=nil {
		resp = Response{"20001","could not found user",err}
		c.Data["json"] = resp
		return
	}
	var userInfoVO responses.UserInfoVO
    userInfoVO.UserName=user.User_name
    userInfoVO.SubjectName=user.Subject_name



//--------------------------------------------------

	data := make(map[string]interface{})
	data["userInfo"] =userInfoVO
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp
}



/**
 8.教师监控页面 （标准差没加）
 */
func (c *SupervisorApiController) TeacherMonitoring() {
	defer c.ServeJSON()
	var requestBody requests.TeacherMonitoring
	var resp  Response
	var err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	//supervisorId := requestBody.SupervisorId
	  questionId := requestBody.QuestionId

	//----------------------------------------------------
   paperDistributions :=make([]models.PaperDistribution ,0)
	err = models.FindPaperDistributionByQuestionId(&paperDistributions, questionId)
	if err!=nil {
		resp = Response{"20021","FindPaperDistributionByQuestionId  fail",err}
		c.Data["json"] = resp
		return
	}
	teacherMonitoringList :=  make([]responses.TeacherMonitoringVO,len(paperDistributions))
	for i :=0 ;i<len(paperDistributions);i++ {
		//教师id
		userId:= paperDistributions[i].User_id
		teacherMonitoringList [i].UserId=userId
       //分配试卷数量
		testDistributionNumber:=paperDistributions[i].Test_distribution_number
		teacherMonitoringList [i].TestDistributionNumber= testDistributionNumber
		//testDistributionNumberString:=strconv.FormatInt(testDistributionNumber,10)
		//testDistributionNumberFloat,_:=strconv.ParseFloat(testDistributionNumberString,64)

		finishCount ,err1:= models.CountFinishTestNumberByUserId(userId,questionId)
		if err1!=nil {
			resp = Response{"20022","CountFinishTestNumberByUserId  fail",err}
			c.Data["json"] = resp
			return
		}
		teacherMonitoringList[i].TestSuccessNumber=finishCount
		finishCountString:=strconv.FormatInt(finishCount,10)
		finishCountFloat,_:=strconv.ParseFloat(finishCountString,64)

		remainingTestNumber,err1 := models.CountRemainingTestNumberByUserId(questionId,userId)
		if err1!=nil {
			resp = Response{"20023","CountRemainingTestNumberByUserId  fail",err}
			c.Data["json"] = resp
			return
		}
		teacherMonitoringList[i].TestRemainingNumber=remainingTestNumber

		failCount,err1:= models.CountFailTestNumberByUserId(userId,questionId)
		if err1!=nil {
			resp = Response{"20024","CountFailTestNumberByUserId  fail",err}
			c.Data["json"] = resp
			return
		}
		teacherMonitoringList[i].TestProblemNumber=failCount
		failCountString:=strconv.FormatInt(finishCount,10)
		failCountFloat,_:=strconv.ParseFloat(failCountString,64)


		user:=models.User{User_id: userId}
		err = user.GetUser(userId)
		if err!=nil {
			resp = Response{"20001","could not found user",err}
			c.Data["json"] = resp
			return
		}
		teacherMonitoringList[i].UserName=user.User_name
  		onlineTime := user.Online_time
		teacherMonitoringList[i].OnlineTime=onlineTime

		var markingSpeed float64 =0
		if onlineTime!=0 {
			markingSpeed =finishCountFloat/onlineTime
		}
		teacherMonitoringList[i].MarkingSpeed=markingSpeed

		var averageScore float64 =0
		if finishCount!=0 {
			sum,err1:=models.SumFinishScore(userId,questionId)
			if err1!=nil {
				resp = Response{"20025","SumFinishScore  fail",err}
				c.Data["json"] = resp
				return
			}
			averageScore=sum/finishCountFloat
		}
		teacherMonitoringList[i].AverageScore=averageScore

		var validity  float64=0
		if (finishCountFloat+failCountFloat)!=0 {
			validity =finishCountFloat/(finishCountFloat+failCountFloat)
		}
		teacherMonitoringList[i].Validity=validity

		selfTestCount ,err1 := models.CountSelfScore(userId,questionId)
		if err1!=nil {
			resp = Response{"20026","CountSelfScore  fail",err}
			c.Data["json"] = resp
			return
		}
		selfTestCountString:=strconv.FormatInt(selfTestCount,10)
		selfTestCountFloat,_:=strconv.ParseFloat(selfTestCountString,64)

		var selfScoreRate float64=0
		if finishCount!=0 {
			selfScoreRate= selfTestCountFloat/finishCountFloat
		}
		teacherMonitoringList[i].EvaluationIndex=selfScoreRate

		//标准差不会
		teacherMonitoringList[i].StandardDeviation=0

	}


//--------------------------------------------------

	data := make(map[string]interface{})

	data["teacherMonitoringList"] =teacherMonitoringList
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp
}




/**
11.分数分布表
*/
func (c *SupervisorApiController) ScoreDistribution() {
	defer c.ServeJSON()
	var requestBody requests.ScoreDistribution
	var resp Response
	var  err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	//supervisorId := requestBody.SupervisorId
	questionId :=requestBody.QuestionId

	//----------------------------------------------------
//求大题满分
    topic :=models.Topic{Question_id: questionId}
	err = topic.GetTopic(questionId)
	if err!=nil {
		resp = Response{"20002","could not find topic",err}
		c.Data["json"] = resp
		return
	}
	questionScore:=topic.Question_score
//求该大题的已批改试卷表
	scoreRecordList := make([]models.ScoreRecord,0)
    err=models.FindFinishScoreRecordListByQuestionId(&scoreRecordList,questionId)
	if err!=nil {
		resp = Response{"20003","FindFinishScoreRecordListByQuestionId err",err}
		c.Data["json"] = resp
		return

	}
//该题已批改试卷总数
	count :=len(scoreRecordList)
	countString:=strconv.FormatInt(int64(count),10)
	countFloat,_:=strconv.ParseFloat(countString,64)
//标准的输出数据
scoreDistributionList := make([]responses.ScoreDistributionVO,questionScore+1)
//统计分数
var i int64=0
for  ;i<=questionScore;i++{
    scoreDistributionList[i].Score=i
	score, err := models.CountTestByScore(questionId, i)
	if err!=nil {
		resp = Response{"20004","CountTestByScore err",err}
		c.Data["json"] = resp
		return
	}
	number := score
	numberString:=strconv.FormatInt(number,10)
	numberFloat,_:=strconv.ParseFloat(numberString,64)
	scoreDistributionList[i].Rate=numberFloat/countFloat
	}

	//--------------------------------------------------

	data := make(map[string]interface{})
	data["scoreDistributionList"] =scoreDistributionList
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp

}

/**
12.大题教师选择列表
*/
func (c *SupervisorApiController) TeachersByQuestion() {
	defer c.ServeJSON()
	var requestBody requests.TeachersByQuestion
	var resp Response
	var  err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	//supervisorId := requestBody.SupervisorId
	questionId :=requestBody.QuestionId

	//----------------------------------------------------
	//根据大题求试卷分配表
	paperDistributions :=make([]models.PaperDistribution ,0)
	err = models.FindPaperDistributionByQuestionId(&paperDistributions, questionId)
	if err!=nil{resp = Response{"20005","FindPaperDistributionByQuestionId err",err}
		c.Data["json"] = resp
		return}

	//输出标准	
	teacherVOList := make([]responses.TeacherVO,len(paperDistributions))

	//求教师名和转化输出
  	for i:=0 ;i<len(paperDistributions);i++ {
		userId :=paperDistributions[i].User_id
  	    user:=models.User{User_id: userId}
		err := user.GetUser(userId)
		if err!=nil {
			resp = Response{"20001","could not found user",err}
			c.Data["json"] = resp
			return
		}
		userName :=user.User_name
		teacherVOList[i].UserId=userId
		teacherVOList[i].UserName=userName
	}

	//--------------------------------------------------

	data := make(map[string]interface{})
	data["teacherVOList"] =teacherVOList
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp

}

/**
13.自评监控表
*/
func (c *SupervisorApiController) SelfScore() {
	defer c.ServeJSON()
	var requestBody requests.SelfScore
	var resp Response
	var  err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	//supervisorId := requestBody.SupervisorId
	examinerId :=requestBody.ExaminerId
	//----------------------------------------------------

	//根据userId找到自评卷

	selfScoreRecord :=make([]models.ScoreRecord ,0)
	models.FindSelfScoreRecordByUserId(&selfScoreRecord,examinerId)
	if err!=nil {
		resp = Response{"20006","FindSelfScoreRecordByUserId err",err}
		c.Data["json"] = resp
		return
	}
	//输出标准
	selfScoreRecordVOList := make([]responses.SelfScoreRecordVO,len(selfScoreRecord))

	//求教师名和转化输出
  	for i:=0 ;i<len(selfScoreRecord);i++ {
  		testId :=selfScoreRecord[i].Test_id
		var testScoreRecord models.ScoreRecord
  		models.GetTestScoreRecordByTestIdAndUserId(&testScoreRecord,testId,examinerId)
  	    selfScoreRecordVOList[i].TestId=testId
  	    selfScoreRecordVOList[i].Score=testScoreRecord.Score
  	    selfScoreRecordVOList[i].SelfScore=selfScoreRecord[i].Score

	}

	//--------------------------------------------------

	data := make(map[string]interface{})
	data["selfScoreRecordVOList"] =selfScoreRecordVOList
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp

}
/**
14，平均分监控表
*/
func (c *SupervisorApiController) AverageScore() {
	defer c.ServeJSON()
	var requestBody requests.AverageScore
	var resp Response
	var  err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	//supervisorId := requestBody.SupervisorId
	questionId :=requestBody.QuestionId
//--------------------------------------------
	//根据大题求试卷分配表
	paperDistributions :=make([]models.PaperDistribution ,0)
	err = models.FindPaperDistributionByQuestionId(&paperDistributions, questionId)
	if err!=nil {
		resp = Response{"20007","FindPaperDistributionByQuestionId err",err}
		c.Data["json"] = resp
		return
	}
	//输出标准
	scoreAverageVOList := make([]responses.ScoreAverageVO,len(paperDistributions))

	//求教师名和转化输出
  	for i:=0 ;i<len(paperDistributions);i++ {
		//求userId 和userName
		userId :=paperDistributions[i].User_id
		user:=models.User{User_id: userId}
		err := user.GetUser(userId)
		if err!=nil {
			resp = Response{"20001","could not found user",err}
			c.Data["json"] = resp
			return
		}
		userName :=user.User_name
		scoreAverageVOList[i].UserId=userId
		scoreAverageVOList[i].UserName=userName

		finishCount, err := models.CountFinishTestNumberByUserId(userId, questionId)
		if err!=nil {
			resp = Response{"20008","CountFinishTestNumberByUserId fail",err}
			c.Data["json"] = resp
			return
		}
		finishCountString:=strconv.FormatInt(finishCount,10)
		finishCountFloat,_:=strconv.ParseFloat(finishCountString,64)

		var averageScore float64 =0
		if finishCount!=0 {
			sum, err := models.SumFinishScore(userId, questionId)
			if err!=nil {
				resp = Response{"20009","SumFinishScore fail",err}
				c.Data["json"] = resp
				return
			}
			averageScore=sum/finishCountFloat
		}
		scoreAverageVOList[i].Average=averageScore

	}
	var topic =models.Topic{Question_id: questionId}
	err = topic.GetTopic(questionId)
	if err!=nil {
		resp  = Response{"20000","GetTopicList err ",err}
		c.Data["json"] = resp
		return
	}
	var fullScore =topic.Question_score


	//--------------------------------------------------

	data := make(map[string]interface{})
	data["scoreAverageVOList"] =scoreAverageVOList
	data["fullScore"] =fullScore
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp

}


/**
17，问题卷表
*/
func (c *SupervisorApiController) ProblemTest() {
	defer c.ServeJSON()
	var requestBody requests.ProblemTest
	var resp Response
	var  err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	//supervisorId := requestBody.SupervisorId
	questionId :=requestBody.QuestionId
//------------------------------------------------


	//根据大题号找到问题卷
	problemUnderCorrectedPaper :=make([]models.UnderCorrectedPaper ,0)
	models.FindProblemUnderCorrectedPaperByQuestionId(&problemUnderCorrectedPaper,questionId)
	if err!=nil {
		resp = Response{"20010","FindProblemUnderCorrectedPaperByQuestionId  fail",err}
		c.Data["json"] = resp
		return
	}

	//问题卷的数量
	var  count =len(problemUnderCorrectedPaper)
	//输出标准
	ProblemUnderCorrectedPaperVOList := make([]responses.ProblemUnderCorrectedPaperVO,count)


	//求阅卷老师名和转化输出
  	for i:=0 ;i<len(problemUnderCorrectedPaper);i++ {
		//存testId
  		ProblemUnderCorrectedPaperVOList[i].TestId=problemUnderCorrectedPaper[i].Test_id
		//存userId  userName
		userId :=problemUnderCorrectedPaper[i].User_id
		user:=models.User{User_id: userId}
		err := user.GetUser(userId)
		if err!=nil {
			resp = Response{"20001","could not found user",err}
			c.Data["json"] = resp
			return
		}
		userName :=user.User_name
		ProblemUnderCorrectedPaperVOList[i].ExaminerId=userId
		ProblemUnderCorrectedPaperVOList[i].ExaminerName=userName
		//存问题类型
		ProblemUnderCorrectedPaperVOList[i].ProblemType=problemUnderCorrectedPaper[i].Problem_type
		ProblemUnderCorrectedPaperVOList[i].ProblemMes=problemUnderCorrectedPaper[i].Problem_message

	}
	

	//--------------------------------------------------

	data := make(map[string]interface{})
	data["ProblemUnderCorrectedPaperVOList"] =ProblemUnderCorrectedPaperVOList
	data["count"] =count
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp

}

/**
18，仲裁卷表
*/
func (c *SupervisorApiController) ArbitramentTest() {
	defer c.ServeJSON()
	var requestBody requests.ArbitramentTest
	var resp Response
	var  err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	//supervisorId := requestBody.SupervisorId
	questionId :=requestBody.QuestionId
	//------------------------------------------------

	//根据大题号找到仲裁卷
	arbitramentUnderCorrectedPaper :=make([]models.UnderCorrectedPaper ,0)
	err = models.FindArbitramentUnderCorrectedPaperByQuestionId(&arbitramentUnderCorrectedPaper, questionId)
	if err!=nil {
		resp = Response{"20011","FindArbitramentUnderCorrectedPaperByQuestionId  fail",err}
		c.Data["json"] = resp
		return
	}
	//输出标准
	arbitramentTestVOList := make([]responses.ArbitramentTestVO,0)

	var count = len(arbitramentUnderCorrectedPaper)
	//求阅卷老师名和转化输出
  	for i:=0 ;i<len(arbitramentUnderCorrectedPaper);i++ {
		//存testId
		var testId = arbitramentUnderCorrectedPaper[i].Test_id
  		arbitramentTestVOList[i].TestId=testId

		//查试卷
		var testPaper models.TestPaper
  		testPaper.Test_id=testId

  		testPaper.GetTestPaper(testId)
		//查存试卷第一次评分人id
		var examinerFirstId = testPaper.Examiner_first_id
  		arbitramentTestVOList[i].ExaminerFirstId=examinerFirstId
		//查第一次评分人
		firstExaminer:=models.User{User_id: examinerFirstId}
		err := firstExaminer.GetUser(examinerFirstId)
		if err!=nil {
			resp = Response{"20001","could not found user",err}
			c.Data["json"] = resp
			return
		}
		//查第一次评分人姓名
		examinerFirstName :=firstExaminer.User_name
        //存试卷第一次评分人姓名和分数
		arbitramentTestVOList[i].ExaminerFirstName=examinerFirstName
		arbitramentTestVOList[i].ExaminerFirstScore=testPaper.Examiner_first_score
		//查存试卷第二次评分人id
		var examinerSecondId = testPaper.Examiner_second_id
		arbitramentTestVOList[i].ExaminerSecondId=examinerSecondId
		//查第二次试卷评分人
		secondExaminer:=models.User{User_id: examinerSecondId}
		err = secondExaminer.GetUser(examinerSecondId)
		if err!=nil {
			resp = Response{"20001","could not found user",err}
			c.Data["json"] = resp
			return
		}
		//查第二次评分人姓名
		secondExaminerName :=secondExaminer.User_name
		//存第一次评分人姓名和分数
		arbitramentTestVOList[i].ExaminerSecondName=secondExaminerName
		arbitramentTestVOList[i].ExaminerSecondScore=testPaper.Examiner_second_score
		//查存实际误差
		arbitramentTestVOList[i].PracticeError=testPaper.Pratice_error
		//查存标准误差
		var topic  models.Topic
		topic.GetTopic(questionId)
		arbitramentTestVOList[i].StandardError=topic.Standard_error

	}
	//查存该题满分
	var topic =models.Topic{Question_id: questionId}
  	topic.GetTopic(questionId)
	var fullScore =topic.Question_score


	//--------------------------------------------------

	data := make(map[string]interface{})
	data["arbitramentTestVOList"] =arbitramentTestVOList
	data["count"] =count
	data["fullScore"] =fullScore
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp

}

/**
15.总体进度（平均分没加）
*/
func (c *SupervisorApiController) ScoreProgress() {
	defer c.ServeJSON()
	var requestBody requests.ArbitramentTest
	var resp Response
	var  err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
	//supervisorId := requestBody.SupervisorId

	//----------------------------------------------------
	//获取大题列表
	topics :=make([]models.Topic ,0)
	err = models.GetTopicList(&topics)
	if err!=nil {
		resp  = Response{"20000","GetTopicList err ",err}
		c.Data["json"] = resp
		return
	}
	//确定输出标准
	scoreProgressVOList :=  make([]responses.ScoreProgressVO,len(topics))

	for i :=0 ;i<len(topics);i++ {
		//获取大题id
		questionId := topics[i].Question_id
		scoreProgressVOList[i].QuestionId = questionId
		//获取大题名
		questionName := topics[i].Question_name
		scoreProgressVOList[i].QuestionName = questionName
		//获取 任务总量
		importNumber := topics[i].Import_number
		scoreProgressVOList[i].ImportNumber = importNumber

		//出成绩量
		finishNumber,err1 := models.CountFinishScoreNumberByQuestionId(questionId)
		if err1!=nil {
			resp = Response{"20013","CountFinishScoreNumberByQuestionId  fail",err}
			c.Data["json"] = resp
			return
		}
		scoreProgressVOList[i].FinishNumber = finishNumber
		//出成绩率
		finishNumberString := strconv.FormatInt(finishNumber, 10)
		finishNumberFloat, _ := strconv.ParseFloat(finishNumberString, 64)
		importNumberString := strconv.FormatInt(importNumber, 10)
		importNumberFloat, _ := strconv.ParseFloat(importNumberString, 64)
		var finishRate float64 = 0
		if importNumberFloat != 0 {
			finishRate = finishNumberFloat / importNumberFloat
		}
		scoreProgressVOList[i].FinishRate = finishRate
		//未出成绩量
		unfinishedNumberFloat := importNumberFloat - finishNumberFloat
		scoreProgressVOList[i].UnfinishedNumber = unfinishedNumberFloat
		//未出成绩率
		var unfinishedRate float64 = 0
		if (importNumberFloat != 0) {
			unfinishedRate = unfinishedNumberFloat / importNumberFloat
		}
		scoreProgressVOList[i].UnfinishedRate = unfinishedRate
		//是否全部完成
		var isAllFinished int64
		if unfinishedNumberFloat != 0 {
			isAllFinished = 0
		} else {
			isAllFinished = 1
		}
		scoreProgressVOList[i].IsAllFinished = isAllFinished
		//--------------------------------------------------
		//一次评卷完成数
		firstScoreNumber ,err1:= models.CountFirstScoreNumberByQuestionId(questionId)
		if err1!=nil {
			resp = Response{"20014","CountFirstScoreNumberByQuestionId  fail",err}
			c.Data["json"] = resp
			return
		}
		scoreProgressVOList[i].FirstFinishedNumber = firstScoreNumber
		//一次评卷完成率
		firstScoreNumberString := strconv.FormatInt(firstScoreNumber, 10)
		firstScoreNumberFloat, _ := strconv.ParseFloat(firstScoreNumberString, 64)
		var firstScoreRate float64 = 0
		if (importNumberFloat != 0) {
			firstScoreRate = firstScoreNumberFloat / importNumberFloat
		}
		scoreProgressVOList[i].FirstFinishedRate = firstScoreRate
		//未出第一次成绩量
		firstUnfinishedNumber := importNumberFloat - firstScoreNumberFloat
		scoreProgressVOList[i].FirstUnfinishedNumber = firstUnfinishedNumber
		//第一次未出成绩率
		var firstUnfinishedRate float64 = 0
		if (importNumberFloat != 0) {
			firstUnfinishedRate = firstUnfinishedNumber / importNumberFloat
		}
		scoreProgressVOList[i].FirstUnfinishedRate = firstUnfinishedRate
		//第一次阅卷是否全部完成
		var isFirstFinished int64
		if firstUnfinishedNumber != 0 {
			isAllFinished = 0
		} else {
			isAllFinished = 1
		}
		scoreProgressVOList[i].IsFirstFinished = isFirstFinished

		//-----------------------------------------

		//二次评卷完成数
		secondScoreNumber,err1 := models.CountSecondScoreNumberByQuestionId(questionId)
		if err1!=nil {
			resp = Response{"20015","CountSecondScoreNumberByQuestionId  fail",err}
			c.Data["json"] = resp
			return
		}
		scoreProgressVOList[i].SecondFinishedNumber = secondScoreNumber
		//二次评卷完成率
		secondScoreNumberString := strconv.FormatInt(secondScoreNumber, 10)
		secondScoreNumberFloat, _ := strconv.ParseFloat(secondScoreNumberString, 64)
		var secondScoreRate float64 = 0
		if (importNumberFloat != 0) {
			secondScoreRate = secondScoreNumberFloat / importNumberFloat
		}
		scoreProgressVOList[i].SecondFinishedRate = secondScoreRate

		//未出第二次成绩量
		secondUnfinishedNumber := importNumberFloat - firstScoreNumberFloat
		scoreProgressVOList[i].SecondUnfinishedNumber = secondUnfinishedNumber
		//第二次未出成绩率
		var secondUnfinishedRate float64 = 0
		if (importNumberFloat != 0) {
			secondUnfinishedRate = secondUnfinishedNumber / importNumberFloat
		}
		scoreProgressVOList[i].SecondUnfinishedRate = secondUnfinishedRate
		//第二次阅卷是否全部完成
		var isSecondFinished int64
		if secondUnfinishedNumber != 0 {
			isSecondFinished = 0
		} else {
			isSecondFinished = 1
		}
		scoreProgressVOList[i].IsSecondFinished = isSecondFinished

		//-----------------------------------------

		//三次评卷完成数
		thirdScoreNumber,err1 := models.CountThirdScoreNumberByQuestionId(questionId)
		if err1!=nil {
			resp = Response{"20016","CountThirdScoreNumberByQuestionId  fail",err}
			c.Data["json"] = resp
			return
		}
		scoreProgressVOList[i].ThirdFinishedNumber = thirdScoreNumber
		//三次评卷完成率
		thirdScoreNumberString := strconv.FormatInt(thirdScoreNumber, 10)
		thirdScoreNumberFloat, _ := strconv.ParseFloat(thirdScoreNumberString, 64)
		var thirdScoreRate float64 = 0
		if (importNumberFloat != 0) {
			thirdScoreRate = thirdScoreNumberFloat / importNumberFloat
		}
		scoreProgressVOList[i].ThirdFinishedRate = thirdScoreRate

		//未出第三次成绩量
		thirdUnfinishedNumber := importNumberFloat - thirdScoreNumberFloat
		scoreProgressVOList[i].ThirdUnfinishedNumber = thirdUnfinishedNumber
		//未出成绩率
		var thirdUnfinishedRate float64 = 0
		if (importNumberFloat != 0) {
			thirdUnfinishedRate = thirdUnfinishedNumber / importNumberFloat
		}
		scoreProgressVOList[i].ThirdUnfinishedRate = thirdUnfinishedRate
		//第三次阅卷是否全部完成
		var isThirdFinished int64
		if thirdUnfinishedNumber != 0 {
			isThirdFinished = 0
		} else {
			isThirdFinished = 1
		}
		scoreProgressVOList[i].IsThirdFinished = isThirdFinished

		//-----------------------------------------
		//仲裁卷完成数
		arbitramentFinishNumber,err1 := models.CountArbitramentFinishNumberByQuestionId(questionId)
		if err1!=nil {
			resp = Response{"20017","CountArbitramentFinishNumberByQuestionId  fail",err}
			c.Data["json"] = resp
			return
		}
		scoreProgressVOList[i].ArbitramentFinishedNumber = arbitramentFinishNumber
		//仲裁卷未完成量
		arbitramentUnfinishedNumber,err1 := models.CountArbitramentUnFinishNumberByQuestionId(questionId)
		if err1!=nil {
			resp = Response{"20018","CountArbitramentUnFinishNumberByQuestionId  fail",err}
			c.Data["json"] = resp
			return
		}

		scoreProgressVOList[i].ArbitramentUnfinishedNumber = arbitramentUnfinishedNumber
		//仲裁卷产生量：
		arbitramentNumber := arbitramentFinishNumber+arbitramentUnfinishedNumber
		scoreProgressVOList[i].ArbitramentNumber = arbitramentNumber
		//仲裁卷产生率
		arbitramentNumberString := strconv.FormatInt(arbitramentNumber, 10)
		arbitramentNumberFloat, _ := strconv.ParseFloat(arbitramentNumberString, 64)
		var arbitramentRate float64 = 0
		if (importNumberFloat != 0) {
			arbitramentRate = arbitramentNumberFloat / importNumberFloat
		}
		scoreProgressVOList[i].ArbitramentRate = arbitramentRate

		//仲裁卷完成率
		arbitramentFinishNumberString := strconv.FormatInt(arbitramentFinishNumber, 10)
		arbitramentFinishNumberFloat, _ := strconv.ParseFloat(arbitramentFinishNumberString, 64)
		var arbitramentFinishRate float64 = 0
		if (arbitramentNumberFloat != 0) {
			arbitramentFinishRate = arbitramentFinishNumberFloat / arbitramentNumberFloat
		}
		scoreProgressVOList[i].ArbitramentFinishedRate = arbitramentFinishRate


		//仲裁卷未完成率
		arbitramentUnFinishNumberString := strconv.FormatInt(arbitramentUnfinishedNumber, 10)
		arbitramentUnFinishNumberFloat, _ := strconv.ParseFloat(arbitramentUnFinishNumberString, 64)
		var arbitramentUnfinishedRate float64 = 0
		if (arbitramentNumberFloat != 0) {
			arbitramentUnfinishedRate = arbitramentUnFinishNumberFloat / arbitramentNumberFloat
		}
		scoreProgressVOList[i].ArbitramentUnfinishedRate = arbitramentUnfinishedRate
		//仲裁卷是否全部完成
		var ArbitramentFinished int64
		if arbitramentUnfinishedNumber != 0 {
			ArbitramentFinished = 0
		} else {
			ArbitramentFinished = 1
		}
		scoreProgressVOList[i].IsArbitramentFinished = ArbitramentFinished

		//-----------------------------------------

		//问题卷完成数
		problemFinishNumber,err1 := models.CountProblemFinishNumberByQuestionId(questionId)
		if err1!=nil {
			resp = Response{"20019","CountProblemFinishNumberByQuestionId  fail",err}
			c.Data["json"] = resp
			return
		}
		scoreProgressVOList[i].ProblemFinishedNumber = problemFinishNumber
		//问题卷未完成量
		problemUnfinishedNumber,err1 :=  models.CountProblemUnFinishNumberByQuestionId(questionId)
		if err1!=nil {
			resp = Response{"20020","CountProblemUnFinishNumberByQuestionId  fail",err}
			c.Data["json"] = resp
			return
		}

		scoreProgressVOList[i].ProblemUnfinishedNumber = problemUnfinishedNumber

		//问题卷产生量：
		problemNumber := problemFinishNumber+problemUnfinishedNumber
		scoreProgressVOList[i].ProblemNumber = problemNumber

		//问题卷产生率
		problemNumberString := strconv.FormatInt(problemNumber, 10)
		problemNumberFloat, _ := strconv.ParseFloat(problemNumberString, 64)
		var problemRate float64 = 0
		if (importNumberFloat != 0) {
			problemRate = problemNumberFloat / importNumberFloat
		}
		scoreProgressVOList[i].ProblemRate = problemRate


		//问题卷完成率
		problemFinishedNumberString := strconv.FormatInt(problemFinishNumber, 10)
		problemFinishNumberFloat, _ := strconv.ParseFloat(problemFinishedNumberString, 64)
		var problemFinishRate float64 = 0
		if (problemNumberFloat != 0) {
			problemFinishRate = problemFinishNumberFloat / problemNumberFloat
		}
		scoreProgressVOList[i].ProblemFinishedRate = problemFinishRate


		//问题卷未完成率
		problemUnfinishedNumberrString := strconv.FormatInt(problemUnfinishedNumber, 10)
		problemUnfinishedNumberFloat, _ := strconv.ParseFloat(problemUnfinishedNumberrString, 64)
		var problemUnfinishedRate float64 = 0
		if (problemNumberFloat != 0) {
			problemUnfinishedRate = problemUnfinishedNumberFloat / problemNumberFloat
		}
		scoreProgressVOList[i].ProblemUnfinishedRate = problemUnfinishedRate
		//问题卷是否全部完成
		var IsProblemFinished int64
		if problemUnfinishedNumber != 0 {
			IsProblemFinished = 0
		} else {
			IsProblemFinished = 1
		}
		scoreProgressVOList[i].IsProblemFinished = IsProblemFinished
	}
	//--------------------------------------------------

	data := make(map[string]interface{})

	data["scoreProgressVOList"] =scoreProgressVOList
	resp = Response{"10000", "OK", data}
	c.Data["json"] = resp
}
/**
19.阅卷组长批改试卷
 */
func (c *SupervisorApiController) SupervisorPoint() {
	defer c.ServeJSON()
	var requestBody requests.SupervisorPoint
	var resp Response
	var  err error

	err=json.Unmarshal(c.Ctx.Input.RequestBody, &requestBody)
	if err!=nil {
		resp = Response{"10001","cannot unmarshal",err}
		c.Data["json"] = resp
		return
	}
    supervisorId := requestBody.SupervisorId
    testId := requestBody.TestId
    scoreStr:= requestBody.Scores
    testDetailIdStr:=requestBody.TestDetailIds
	testDetailIds := strings.Split(testDetailIdStr, "-")
	scores := strings.Split(scoreStr, "-")

	//---------------------------------------------------------------------------------------
	//创建试卷小题详情

	var test models.TestPaper

	var sum int64
    //给试卷详情表打分
	for i := 0; i < len(testDetailIds); i++ {
		//取出小题试卷id,和小题分数
		var testInfo models.TestPaperInfo
		testDetailIdString:=testDetailIds[i]
		testDetailId, _ := strconv.ParseInt(testDetailIdString, 10, 64)
		scoreString:=scores[i]
		score, _ := strconv.ParseInt(scoreString, 10, 64)
		//查试卷小题
		err := testInfo.GetTestPaperInfo(testDetailId)
		if err != nil {
			resp := Response{"10008", "get testPaper fail", err}
			c.Data["json"] = resp
			return
		}
		//修改试卷详情表
		testInfo.Leader_id=supervisorId
        testInfo.Leader_score=score
        testInfo.Final_score=score
		err = testInfo.Update()
		if err != nil {
			resp := Response{"10009", "update testPaper fail", err}
			c.Data["json"] = resp
			return
		}
		sum += score
	}
	//给试卷表打分
	err = test.GetTestPaper(testId)
	if err != nil || test.Test_id == 0 {
		resp := Response{"10002", "get test paper fail", err}
		c.Data["json"] = resp
		return
	}
	test.Leader_id = supervisorId
	test.Leader_score = sum
	test.Final_score = sum
	err = test.Update()
	if err != nil {
		resp := Response{"10007", "update test fail", err}
		c.Data["json"] = resp
		return
	}
	//删除试卷待批改表 ，增加试卷记录表
	var record models.ScoreRecord
	var underTest models.UnderCorrectedPaper
	err = models.GetUnderCorrectedPaperByUserIdAndTestId(&underTest, supervisorId, testId)
	if err!=nil {
		resp = Response{"20012","GetUnderCorrectedPaperByUserIdAndTestId  fail",err}
		c.Data["json"] = resp
		return
	}
	record.Score = sum
	record.Test_id = testId
	record.Test_record_type = underTest.Test_question_type
	record.User_id = supervisorId
	record.Question_id=underTest.Question_id
	record.Problem_type=underTest.Problem_type

	err = record.Save()
	if err!=nil {
		resp = Response{"20013","Save  fail",err}
		c.Data["json"] = resp
		return
	}
	err = underTest.Delete()
	if err!=nil {
		resp = Response{"20014","Delete  fail",err}
		c.Data["json"] = resp
		return
	}
	//----------------------------------------
	resp = Response{"10000", "OK", nil}
	c.Data["json"] = resp
}

//16标准差