import React, {Component} from "react";
import DocumentTitle from "react-document-title";
import {Progress, Select, Table} from "antd";
import * as Settings from "../../../../Setting";
import "./index.less";
import group from "../../../../api/group";
import Manage from "../../../../api/manage";
const {Option} = Select;
export default class index extends Component {

  componentDidMount() {
    this.questionList();
  }

  // 选择区
  state = {
    questionList: [],
    tableData: [],
    fullScore: undefined,
    averageScore: undefined,
    subjectList: [],

  }

  questionList = () => {
    Manage.subjectList().then((res) => {
      this.setState({subjectList: res.data.data.subjectVOList});
    })
      .catch((e) => {
        Settings.showMessage("error", e);
      });
  }
  // 题目选择区
  selectBox = () => {
    let selectList;
    if (this.state.questionList.length !== 0) {
      selectList = this.state.questionList.map((item, i) => {
        return <Option key={i} value={item.QuestionName} label={item.QuestionName}>{item.QuestionName}</Option>;
      });
    } else {
      return null;
    }
    return selectList;
  }
  select = (e) => {
    let index;
    for (let i = 0; i < this.state.questionList.length; i++) {
      if (this.state.questionList[i].QuestionName === e) {
        index = i;
      }
    }
    this.tableData(this.state.questionList[index].QuestionId);
  }
  tableData = (questionId) => {
    group.averageMonitor({supervisorId: "2", questionId: questionId})
      .then((res) => {
        if (res.data.status === "10000") {
          let tableData = [];
          for (let i = 0; i < res.data.data.scoreAverageVOList.length; i++) {
            let item = res.data.data.scoreAverageVOList[i];
            tableData.push({
              UserName: item.UserName,
              Average: item.Average,
            });
          }
          this.setState({
            tableData,
            fullScore: res.data.data.fullScore,
            averageScore: res.data.data.questionAverageScore,
          });
        }
      })
      .catch((e) => {
        Settings.showMessage("error", e);
      });
  }
  columns = [
    {
      title: "教师",
      width: 150,
      dataIndex: "UserName",
    },
    {
      title: "平均分",
      width: 180,
      dataIndex: "Average",
    },
  ]
  progressItem = () => {
    let progressItem = [];
    if (this.state.tableData.length !== 0) {
      progressItem = this.state.tableData.map((item, index) => {
        return <div className="progress-item" key={index}>
          <span>{item.UserName}&nbsp;&nbsp;</span>
          <Progress percent={item.Average / this.state.fullScore * 100} showInfo={false} />
          &nbsp;&nbsp;<span>{item.Average}</span>
        </div>;
      });
    } else {
      return null;
    }
    return progressItem;
  }

  onSelectsub = (e) => {
    group.questionList({subjectName: e})
      .then((res) => {
        if (res.data.status === "10000") {
          this.setState({
            questionList: res.data.data.questionsList,
          });
          if (res.data.data.questionsList.length > 0) {this.tableData(res.data.data.questionsList[0].QuestionId);}
        }
      });
  }

  selectSubject = () => {
    return this.state.subjectList.map((item, i) => {
      return <Option key={i} value={item.SubjectName} label={item.SubjectName}>{item.SubjectName}</Option>;
    });
  }
  render() {
    return (
      <DocumentTitle title="阅卷系统-平均分监控">
        <div className="average-monitor-page" data-component="average-monitor-page">
          <div className="search-container">
            <div className="question-select">
                题目选择：<Select
                showSearch
                style={{width: 120, marginRight: 70}}
                optionFilterProp="label"
                onSelect={(e) => {this.select(e);}}
                filterOption={(input, option) =>
                  option.label.indexOf(input) >= 0
                }
                filterSort={(optionA, optionB) =>
                  optionA.label.localeCompare(optionB.label)
                }
                placeholder={this.state.questionList.length > 0 ? this.state.questionList[0].QuestionName : null}
                defaultValue={this.state.questionList.length > 0 ? this.state.questionList[0].QuestionName : null}
              >
                {
                  this.selectBox()
                }

              </Select>
            </div>
            科目选择：<Select
              style={{width: 120}}
              optionFilterProp="label"
              onSelect={(e) => {this.onSelectsub(e);}}
              filterOption={(input, option) =>
                option.label.indexOf(input) >= 0
              }
              filterSort={(optionA, optionB) =>
                optionA.label.localeCompare(optionB.label)
              }>
              {
                this.selectSubject()
              }
            </Select>
            <div className="question-score">
                满分：{this.state.fullScore}
            </div>
            <div className="question-score">
                平均分：{this.state.averageScore}
            </div>
          </div>
          <div className="display-container">
            <div className="display-table">
              <Table
                pagination={{position: ["bottomCenter"]}}
                columns={this.columns}
                dataSource={this.state.tableData}
              />
            </div>
            <div className="progress">
              <div className="progress-header">
                <span>教师</span><span>平均分</span>
              </div>
              {
                this.progressItem()
              }
              <div className="progress-footer">
                <span>满分：{this.state.fullScore}</span>
              </div>
            </div>
          </div>

        </div>
      </DocumentTitle>
    );
  }

}
