import React, { useState } from "react";
import {
  Form,
  Input,
  Select,
  Radio,
  Space,
  Typography,
  Checkbox,
  Spin,
  Button,
  message,
} from "antd";

const { Title } = Typography;

const QuestionCreateByHand = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [contentLength, setContentLength] = useState(0);

  // 验证答案：单选题必须选一个，多选题至少选两个
  const answerValidator = (_, value) => {
    const questionType = form.getFieldValue("type");
    if (questionType === "multiple") {
      return Array.isArray(value) && value.length >= 2
        ? Promise.resolve()
        : Promise.reject(new Error("多选题请至少选择两个答案"));
    }
    if (typeof value === "string" && value) {
      return Promise.resolve();
    }
    return Promise.reject(new Error("单选题请选中一个答案"));
  };

  // 构造提交体并发起 POST 请求
  const handleSubmit = async (values) => {
    const payload = {
      type: values.type === "multiple" ? 2 : 1,
      title: values.title,
      language: values.language,
      answers: ["A", "B", "C", "D"].map(
        (k) => `${k}: ${values[`option${k}`] || ""}`
      ),
      rights: Array.isArray(values.answer) ? values.answer : [values.answer],
    };

    try {
      setLoading(true);
      const res = await fetch("/api/questions/CreateByHand", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const json = await res.json();
      if (json.code !== 0) throw new Error(json.msg || "创建失败");
      message.success("创建成功");
      form.resetFields();
      setContentLength(0);
    } catch (err) {
      console.error("创建错误", err);
      message.error("创建失败：" + err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Spin spinning={loading}>
      <div
        style={{
          margin: "0 auto",
          padding: 24,
          backgroundColor: "white",
          borderRadius: 8,
        }}
      >
        <Title
          level={4}
          style={{ color: "#1777FF", marginBottom: 10, fontSize: 20 }}
        >
          创建题目
        </Title>

        <Form form={form} layout="horizontal" onFinish={handleSubmit}>
          <Space size={24} style={{ width: "100%" }}>
            <Form.Item
              label="题型"
              name="type"
              rules={[{ required: true, message: "请选择题型" }]}
            >
              <Select style={{ width: 150 }}>
                <Select.Option value="single">单选题</Select.Option>
                <Select.Option value="multiple">多选题</Select.Option>
                <Select.Option value="programming">编程题</Select.Option>
              </Select>
            </Form.Item>

            <Form.Item
              label="语言"
              name="language"
              rules={[{ required: true, message: "请选择语言" }]}
            >
              <Select style={{ width: 320 }}>
                <Select.Option value="go">Go 语言</Select.Option>
                <Select.Option value="java">Java 语言</Select.Option>
                <Select.Option value="python">Python 语言</Select.Option>
                <Select.Option value="javascript">
                  JavaScript 语言
                </Select.Option>
                <Select.Option value="c++">C++ 语言</Select.Option>
              </Select>
            </Form.Item>
          </Space>

          <Form.Item
            label="标题"
            name="title"
            rules={[{ required: true, message: "请输入标题" }]}
          >
            <Input placeholder="请输入标题" style={{ flex: 1 }} />
          </Form.Item>

          <Form.Item label="内容" name="content">
            <Input.TextArea
              minrows={4}
              maxrows={8}
              rows={7}
              onChange={(e) => setContentLength(e.target.value.length)}
              style={{ flex: 1 }}
            />
            <span style={{ fontSize: 13, color: "black", marginLeft: 8 }}>
              字符数：{contentLength}
            </span>
          </Form.Item>

          {["A", "B", "C", "D"].map((opt) => (
            <Form.Item
              key={opt}
              label={`选项${opt}`}
              name={`option${opt}`}
              rules={[{ required: true, message: `请输入选项${opt}` }]}
            >
              <Input placeholder={`选项${opt}`} />
            </Form.Item>
          ))}

          {/* 答案选择 & 动态切换 & 校验 */}
          <Form.Item shouldUpdate noStyle>
            {({ getFieldValue }) => (
              <Form.Item
                label="答案"
                name="answer"
                rules={[{ validator: answerValidator }]}
              >
                {getFieldValue("type") === "multiple" ? (
                  <Checkbox.Group>
                    <Space>
                      {["A", "B", "C", "D"].map((o) => (
                        <Checkbox key={o} value={o}>
                          {o}
                        </Checkbox>
                      ))}
                    </Space>
                  </Checkbox.Group>
                ) : (
                  <Radio.Group>
                    <Space>
                      {["A", "B", "C", "D"].map((o) => (
                        <Radio key={o} value={o}>
                          {o}
                        </Radio>
                      ))}
                    </Space>
                  </Radio.Group>
                )}
              </Form.Item>
            )}
          </Form.Item>

          <Form.Item
            wrapperCol={{ offset: 6, span: 18 }}
            style={{ textAlign: "right", marginRight: 170 }}
          >
            <Space size={24}>
              <Button onClick={() => window.history.back()}>取消</Button>
              <Button type="primary" htmlType="submit">
                确认
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </div>
    </Spin>
  );
};

export default QuestionCreateByHand;
