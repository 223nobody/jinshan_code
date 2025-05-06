import React, { useState } from "react";
import {
  Form,
  Select,
  Radio,
  InputNumber,
  Button,
  Space,
  Typography,
  Spin,
  Card,
  Input,
  Checkbox,
} from "antd";
import { message } from "antd";
import { FontSizeOutlined, PlusOutlined } from "@ant-design/icons";

const { Title } = Typography;

// 修改预览区域和交互逻辑
const QuestionCreateByAI = () => {
  // 新增状态管理
  const [selectedQuestions, setSelectedQuestions] = useState(new Set());
  const [form] = Form.useForm(); // Define the form instance
  const [previewQuestions, setPreviewQuestions] = useState([]); // Define previewQuestions state
  const [loading, setLoading] = useState(false); // Define loading state

  // 优化选中逻辑
  const handleSelect = (index) => {
    setSelectedQuestions((prev) => {
      const newSet = new Set(prev);
      newSet.has(index) ? newSet.delete(index) : newSet.add(index);
      return newSet;
    });
  };

  // 新增添加到题库功能
  // 修改handleAddToBank函数
  const handleAddToBank = async () => {
    try {
      // 验证表单数据
      const values = await form.validateFields();
      if (!values.languages || values.languages.length === 0) {
        return message.error("请至少选择一种编程语言");
      }

      // 构建符合接口要求的请求体
      const payload = Array.from(selectedQuestions).map((i) => {
        const question = previewQuestions[i];
        return {
          type: { single: 1, multiple: 2, programming: 3 }[question.type] || 1,
          title: question.title,
          language: values.languages[0],
          answers: question.answers.map((a) => a.split("：")[1]?.trim() || a),
          rights: question.rights.map((r) => r[0]), // 取第一个字符
        };
      });

      // 发送请求
      const response = await fetch("/api/questions/batch-insert", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      // 处理响应
      const result = await response.json();
      if (response.status !== 200) {
        throw new Error(result.msg || `HTTP错误 ${response.status}`);
      }

      if (result.code !== 0) {
        throw new Error(result.msg || "服务端返回错误");
      }

      message.success(`成功添加${payload.length}题`);
      setSelectedQuestions(new Set());
    } catch (error) {
      console.error("添加失败详情:", {
        error,
        time: new Date().toISOString(),
      });
      message.error(`添加失败: ${error.message}`);
    }
  };
  const handleGenerate = async (values) => {
    setLoading(true);
    try {
      // 生成指定数量的题目
      const requests = Array.from({ length: values.count }).map(() =>
        fetch("/api/questions/CreateByAI", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            model: "", // 根据实际模型填写
            language: values.languages[0], // 取第一个选中的语言
            type: values.type === "multiple" ? 2 : 1, // 类型转换
            keyword: values.keyword,
          }),
        })
      );

      // 处理并发请求
      const responses = await Promise.allSettled(requests);

      // 过滤成功的响应
      const validResponses = responses
        .filter((res) => res.status === "fulfilled" && res.value.ok)
        .map((res) => res.value);

      // 提取题目数据
      const questions = await Promise.all(
        validResponses.map((res) =>
          res.json().then((data) => ({
            title: data.aiRes.title,
            type: values.type, // 保留前端类型
            answers: data.aiRes.answers,
            rights: data.aiRes.rights,
          }))
        )
      );

      // 更新预览区域
      setPreviewQuestions(questions);

      // 显示生成结果提示
      if (questions.length < values.count) {
        message.warning(`成功生成${questions.length}/${values.count}题`);
      } else {
        message.success(`已生成${values.count}道题目`);
      }
    } catch (error) {
      message.error("生成失败：" + error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: "flex", height: "100%" }}>
      <Card style={{ width: 450, marginRight: 16 }}>
        {
          <div style={{ display: "flex", alignItems: "center" }}>
            <img
              src="/image ai.png"
              alt="logo"
              style={{ width: 50, height: 50 }}
            />
            <span style={{ fontSize: 20, fontWeight: "bold" }}>
              AI 生成试题
            </span>
          </div>
        }
        <Form
          form={form}
          onFinish={handleGenerate}
          initialValues={{ type: "single", count: 3 }}
        >
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(2, 1fr)", // 两列等宽
              gap: 16, // 间距
              alignItems: "flex-end", // 底部对齐
            }}
          >
            {/* 题型 */}
            <Form.Item
              label="题型"
              name="type"
              style={{ marginTop: 20 }}
              rules={[{ required: true }]}
            >
              <Select style={{ width: "100%" }}>
                <Select.Option value="single">单选题</Select.Option>
                <Select.Option value="multiple">多选题</Select.Option>
                <Select.Option value="programming">编程题</Select.Option>
              </Select>
            </Form.Item>

            {/* 题目数量 */}
            <Form.Item
              label="题目数量"
              name="count"
              style={{ marginTop: 20 }}
              rules={[{ required: true }]}
            >
              <InputNumber
                min={1}
                max={10}
                style={{ width: "100%" }}
                formatter={(value) => `${value}题`}
              />
            </Form.Item>
          </div>

          <Form.Item
            name="languages"
            label="语言"
            style={{ marginTop: 0 }}
            rules={[{ required: true, message: "请选择语言" }]}
          >
            <Select
              mode="multiple"
              placeholder="选择语言"
              style={{ width: "100%" }}
            >
              <Select.Option value="go">Go 语言</Select.Option>
              <Select.Option value="java">Java 语言</Select.Option>
              <Select.Option value="python">Python 语言</Select.Option>
              <Select.Option value="javascript">JavaScript 语言</Select.Option>
              <Select.Option value="c++">C++ 语言</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="keyword"
            label="关键词"
            style={{ marginTop: 0 }}
            rules={[{ required: true, message: "请输入关键词" }]}
          >
            <Input placeholder="请输入关键词" style={{ flex: 1 }} />
          </Form.Item>

          {/* 生成按钮 */}
          <Form.Item
            style={{
              marginTop: 36,
              display: "flex",
              justifyContent: "center",
            }}
          >
            <Button
              htmlType="submit"
              style={{
                width: 300,
                height: 40,
                borderRadius: 8,
                backgroundColor: "#1777FF",
                border: "1px solid #1890ff",
                color: "#fff",
              }}
              icon={<PlusOutlined style={{ color: "#fff" }} />}
              loading={loading}
            >
              生成并预览题库
            </Button>
          </Form.Item>
        </Form>
      </Card>

      {/* 右侧预览区域 */}
      <Card style={{ flex: 1, overflow: "auto", background: "#F0F8FF" }}>
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            padding: "8px 16px",
            background: "#fff",
            borderRadius: 4,
            boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
          }}
        >
          <span style={{ color: "#666" }}>
            已选 {selectedQuestions.size} 题
          </span>
          <Button
            type="primary"
            size="small"
            disabled={selectedQuestions.size === 0}
            onClick={handleAddToBank}
            style={{
              width: 120,
              height: 32,
              fontSize: 14,
              borderRadius: 4,
            }}
          >
            添加到题库
          </Button>
        </div>

        {previewQuestions.map((q, idx) => (
          <Card
            key={idx}
            size="small"
            style={{
              margin: "12px 16px",
              background: selectedQuestions.has(idx) ? "#E6F4FF" : "#F3F4F6",
              border: selectedQuestions.has(idx)
                ? "1px solid #1890ff"
                : "1px solid #F3F4F6",
              borderRadius: 6,
              cursor: "pointer",
              transition: "all 0.3s",
            }}
            bodyStyle={{
              padding: "12px 16px",
              display: "flex",
              alignItems: "flex-start",
            }}
            onClick={() => handleSelect(idx)}
          >
            {/* Checkbox 选择器 */}
            <Checkbox
              checked={selectedQuestions.has(idx)}
              style={{
                marginRight: 12,
                marginTop: 2,
                flexShrink: 0,
              }}
            />

            {/* 题目内容容器 */}
            <div style={{ flex: 1, fontFamily: "Arial, sans-serif" }}>
              <p
                style={{
                  fontSize: 16,
                  color: "#1a3353",
                  marginBottom: 8,
                  fontWeight: 500,
                }}
              >
                {idx + 1}. {q.title}
              </p>

              <ul
                style={{
                  margin: "8px 0",
                  paddingLeft: 0,
                  listStyleType: "none",
                }}
              >
                {q.answers.map((opt, i) => (
                  <li
                    key={i}
                    style={{
                      marginBottom: 6,
                      fontSize: 15,
                      color: "#595959",
                    }}
                  >
                    {opt}
                  </li>
                ))}
              </ul>

              {/* 参考答案右对齐 */}
              <div
                style={{
                  textAlign: "right",
                }}
              >
                <span
                  style={{
                    color: "#00B365",
                    fontSize: 15,
                    fontWeight: 500,
                  }}
                >
                  参考答案：{q.rights.join(", ")}
                </span>
              </div>
            </div>
          </Card>
        ))}
      </Card>
    </div>
  );
};

export default QuestionCreateByAI;
