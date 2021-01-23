# EvaluateSystem
这是一个基于企业微信应用的角色评分系统后端接口，是基于gin框架和gorm进行开发！

# 以下是其数据库的大致设计
用户表:user
  
| id | 用户id | int |  
| --- | ---- | ---- |
| title | 用户姓名 | string |
| wx_user_id | 企业微信用户id | string |
| avatar | 头像 | stirng |  
| power | 权限 | tinyint |
| type | 身份 | tinyint |  
| del | 是否删除 | tinyint |
| del_time | 删除时间 | datetime |
| update_time | 更新时间 | datetime |
| add_time | 添加时间 | datetime |

题库表：problem

| id | 题目id | int |
| --- | --- | ----|  
| title | 题目标题 | string |
| content | 题目内容 | json |
| del | 是否删除 | tinyint |
| del_time | 删除时间 | datetime |
| update_time | 更新时间 | datetime |
| add_time | 添加时间 | datetime |

试卷表:paper

| id | ID | int |
| --- | --- | --- |
| problem_id | 题目id | int |
| type | 身份 | tinyint |
| score | 分数 | int |
| del | 是否删除 | tinyint |
| del_time | 删除时间 | datetime |
| update_time | 更新时间 | datetime |
| add_time | 添加时间 | datetime |

提交表单表:order

| id | 表单id | int |  
| --- | --- | ---- |  
| user_id | 用户id | int |  
| admin_id | 审核人id | json |  
| content | 试卷内容 | json |  
| state | 审核状态 | tinyint |  
| date | 时间 | date |  
| score |  总得分 | int | 
| del | 是否删除 | tinyint |
| del _time | 删除时间 | datetime |
| update_time | 更新时间 | datetime |
| add_time | 添加时间 | datetime |

撤销原因表:reason 


| id | 原因id | int |
| --- | ---- | --- |
| title | 撤销原因 | string |
| admin_id | 审核人id | int |
| user_id | 提交人id | int |
| order_id | 表单id | int |
| del | 是否删除 | datetime |
| del_time | 删除时间 | datetime |
| update_time | 更新时间 | datetime |
| add_time | 添加时间 | datetime |

企业微信机构配置表:branch

| id | 机构id | int |  
| --- | ---- | --- |  
| corp_id | corp_id | string |  
| agent_id | agent_id | string |  
| secret | 企业密钥 | string |  
| access_token | access_token | string |
| expire_time | 过期时间 | datetime | 
| del | 是否删除 | tinyint |  
| del_time | 删除时间 | datetime |
| update_time | 更新时间 | datetime | 
| add_time | 添加时间 | datetime |

