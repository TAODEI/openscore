import React, {Component} from "react";
import DocumentTitle from "react-document-title";
import {Select, Table} from "antd";
import * as Settings from "../../../../Setting";
import "./index.less";
import group from "../../../../api/group";

import ReactEcharts from "echarts-for-react";
import Manage from "../../../../api/manage";

const {Option} = Select;
export default class index extends Component {
    supervisorId = "2"
    state = {
      questionList: [],
      columns: [
        {
          title: "分数",
          width: 120,
          dataIndex: "Score",
        },
      ],
      tableData: [
        {Score: "教师"},
      ],
      subjectList: [],
    }

    componentDidMount() {
      this.questionList();
    }

    getOption = () => {
      let X_data = [];
      let Y_data = [];
      for (let i = 1; i < this.state.columns.length; i++) {
        X_data.push(this.state.columns[i].title);
      }
      for (const key in this.state.tableData[0]) {
        if (key !== "Score") {
          Y_data.push(this.state.tableData[0][key]);
        }
      }
      return {
        xAxis: {
          name: "分数",
          data: X_data,
        },
        yAxis: {
          name: "教师（占比）",
        },
        series: [{
          name: "分数",
          type: "bar",
          data: Y_data,
        }],
      };
    };

    questionList = () => {
      Manage.subjectList().then((res) => {
        this.setState({subjectList: res.data.data.subjectVOList});
      })
        .catch((e) => {
          Settings.showMessage("error", e);
        });
    }

    tableData = (questionId) => {
      group.scoreMonitor({supervisorId: "1", questionId: questionId})
        .then((res) => {
          if (res.data.status === "10000") {
            let columns = [{
              title: "分数",
              width: 120,
              dataIndex: "Score",
            }];
            let tableData = [{
              Score: "教师",
            }];

            for (let i = 0; i < res.data.data.scoreDistributionList.length; i++) {
              let item = res.data.data.scoreDistributionList[i];
              tableData[0]["score_" + item.Score] = item.Rate;
              columns.push(
                {
                  title: item.Score + 1,
                  width: 180,
                  dataIndex: "score_" + item.Score,
                }
              );
            }
            this.setState({
              tableData, columns,
            });
          }
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
    onSelectSub = (e) => {
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

    select = (e) => {
      let index;
      for (let i = 0; i < this.state.questionList.length; i++) {
        if (this.state.questionList[i].QuestionName === e) {
          index = i;
        }
      }
      this.tableData(this.state.questionList[index].QuestionId);
    }

    render() {
      return (
        <DocumentTitle title="阅卷系统-分值分布">
          <div className="score-monitor-page" data-component="score-monitor-page">
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

                科目选择：<Select
                  style={{width: 120}}
                  optionFilterProp="label"
                  onSelect={(e) => {this.onSelectSub(e);}}
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
              </div>
            </div>
            <div className="display-container">
              <Table

                pagination={{position: ["bottomCenter"]}}
                columns={this.state.columns}
                dataSource={this.state.tableData}
              />
            </div>
            <ReactEcharts option={this.getOption()} style={{width: 762, height: 300}} />
          </div>
        </DocumentTitle>
      );
    }

}
