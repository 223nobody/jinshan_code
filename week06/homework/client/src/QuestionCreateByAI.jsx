import React, { useState, useRef } from "react";
import {
  Form,
  Select,
  InputNumber,
  Button,
  Typography,
  Card,
  Input,
  Checkbox,
  Spin,
} from "antd";
import { message } from "antd";
import { PlusOutlined, LoadingOutlined } from "@ant-design/icons";
import image_ai from "./images/image_ai.png";
const { Title } = Typography;
const antIcon = <LoadingOutlined style={{ fontSize: 48 }} spin />;

const QuestionCreateByAI = () => {
  const [selectedQuestions, setSelectedQuestions] = useState(new Set());
  const [form] = Form.useForm();
  const scrollContainerRef = useRef(null);
  const [previewQuestions, setPreviewQuestions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [generating, setGenerating] = useState(false);

  // 修复的滚动控制方法
  const handleScrollToTop = () => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollTo({
        top: 0,
        behavior: "smooth",
      });
    }
  };

  const handleScrollToBottom = () => {
    if (scrollContainerRef.current) {
      const { scrollHeight, clientHeight } = scrollContainerRef.current;
      scrollContainerRef.current.scrollTo({
        top: scrollHeight - clientHeight,
        behavior: "smooth",
      });
    }
  };

  // 修复的题目选择逻辑
  const handleSelectInternal = (index) => {
    setSelectedQuestions((prev) => {
      const newSet = new Set(prev);
      newSet.has(index) ? newSet.delete(index) : newSet.add(index);
      return newSet;
    });

    // 自动滚动定位修复
    if (scrollContainerRef.current) {
      const questionElement = document.getElementById(`question-${index}`);
      if (questionElement) {
        const containerTop =
          scrollContainerRef.current.getBoundingClientRect().top;
        const elementTop = questionElement.getBoundingClientRect().top;
        const offset = elementTop - containerTop - 50;

        scrollContainerRef.current.scrollTo({
          top: scrollContainerRef.current.scrollTop + offset,
          behavior: "smooth",
        });
      }
    }
  };

  const handleAddToBank = async () => {
    try {
      const values = await form.validateFields();
      if (!values.languages || values.languages.length === 0) {
        return message.error("请至少选择一种编程语言");
      }

      const payload = Array.from(selectedQuestions).map((i) => {
        const question = previewQuestions[i];
        return {
          type:
            {
              single: 1,
              multiple: 2,
              programming: 3,
            }[question.type] || 1,
          title: question.title,
          language: values.languages[0],
          answers: question.answers.map((a) => {
            // 处理带字母前缀的答案
            const match = a.match(/^[A-Z]:\s*(.+)/);
            return match ? match[1] : a;
          }),
          rights: question.rights,
        };
      });

      const response = await fetch("/api/questions/batch-insert", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

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
      console.error("添加失败详情:", { error, time: new Date().toISOString() });
      message.error(`添加失败: ${error.message}`);
    }
  };

  const handleGenerate = async (values) => {
    setLoading(true);
    setGenerating(true);
    try {
      // 修改为单次请求
      const response = await fetch("/api/questions/CreateByAI", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          model: "",
          language: values.languages[0],
          count: values.count,
          type:
            values.type === "multiple"
              ? 2
              : values.type === "programming"
              ? 3
              : 1,
          keyword: values.keyword,
        }),
      });

      if (!response.ok) throw new Error(`HTTP ${response.status}`);

      const { aiRes } = await response.json();

      // 直接设置题目列表
      setPreviewQuestions(
        aiRes.questions.map((q, idx) => ({
          id: idx + Date.now(), // 生成唯一ID
          title: q.title,
          type: values.type,
          answers: q.answers,
          rights: q.rights,
        }))
      );

      message.success(`成功生成${aiRes.questions.length}题`);
      handleScrollToTop();
    } catch (error) {
      message.error("生成失败：" + error.message);
    } finally {
      setLoading(false);
      setGenerating(false);
    }
  };

  return (
    <div style={{ display: "flex", height: "100vh", padding: 24 }}>
      {/* 左侧表单完全保持不变 */}
      <Card
        style={{
          width: 450,
          marginRight: 16,
          borderRadius: 8,
          boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
        }}
      >
        <div
          style={{ display: "flex", alignItems: "center", marginBottom: 24 }}
        >
          <img
            src={image_ai}
            alt="logo"
            style={{ width: 50, height: 50, marginRight: 12 }}
          />
          <Title level={4} style={{ margin: 0 }}>
            AI 生成试题
          </Title>
        </div>
        <Form
          form={form}
          onFinish={handleGenerate}
          initialValues={{ type: "single", count: 3 }}
        >
          {/* 保持原有表单结构不变 */}
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(2, 1fr)",
              gap: 16,
              marginBottom: 5,
            }}
          >
            <Form.Item label="题型" name="type" rules={[{ required: true }]}>
              <Select>
                <Select.Option value="single">单选题</Select.Option>
                <Select.Option value="multiple">多选题</Select.Option>
                <Select.Option value="programming">编程题</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item
              label="题目数量"
              name="count"
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
            label="编程语言"
            rules={[{ required: true, message: "请选择语言" }]}
          >
            <Select
              mode="multiple"
              placeholder="选择语言"
              optionLabelProp="label"
            >
              <Select.Option value="go" label="Go">
                Go 语言
              </Select.Option>
              <Select.Option value="java" label="Java">
                Java 语言
              </Select.Option>
              <Select.Option value="python" label="Python">
                Python 语言
              </Select.Option>
              <Select.Option value="javascript" label="JS">
                JavaScript 语言
              </Select.Option>
              <Select.Option value="c++" label="C++">
                C++ 语言
              </Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="keyword"
            label="关键词"
            rules={[{ required: true, message: "请输入关键词" }]}
          >
            <Input placeholder="例如：面向对象、数组" />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              block
              icon={<PlusOutlined />}
              loading={loading}
              style={{
                height: 40,
                borderRadius: 8,
                backgroundColor: "#1777FF",
                borderColor: "#1777FF",
              }}
            >
              生成并预览题库
            </Button>
          </Form.Item>
        </Form>
      </Card>

      {/* 右侧预览区域修复滚动 */}
      <Card
        style={{
          flex: 1,
          borderRadius: 8,
          boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
          display: "flex",
          flexDirection: "column",
          overflow: "hidden",
          position: "relative", // 新增定位
        }}
        bodyStyle={{
          padding: 0,
          display: "flex",
          flexDirection: "column",
          height: "100%",
        }}
      >
        <div
          style={{
            padding: "12px 24px",
            backgroundColor: "#fff",
            borderBottom: "1px solid #f0f0f0",
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <span style={{ color: "#666" }}>
            已选 {selectedQuestions.size} 题
          </span>
          <Button
            type="primary"
            onClick={handleAddToBank}
            disabled={selectedQuestions.size === 0}
            style={{
              width: 120,
              backgroundColor: "#1777FF",
              borderColor: "#1777FF",
              color: "#fff",
            }}
          >
            添加到题库
          </Button>
        </div>

        {/* 加载状态覆盖层 */}
        {generating && (
          <div
            style={{
              position: "absolute",
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              backgroundColor: "rgba(255, 255, 255, 0.8)",
              zIndex: 999,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              flexDirection: "column",
            }}
          >
            <Spin indicator={antIcon} />
            <div
              style={{
                marginTop: 16,
                color: "#1777FF",
                fontSize: 16,
                fontWeight: 500,
              }}
            >
              AI正在努力生成题目中...
            </div>
          </div>
        )}

        {/* 修复滚动容器 */}
        <div
          ref={scrollContainerRef}
          style={{
            flex: 1,
            overflowY: "auto",
            padding: "16px 24px",
            height: "calc(100vh - 180px)",
            position: "relative",
            opacity: generating ? 0.5 : 1, // 添加透明度变化
            transition: "opacity 0.3s",
            scrollbarGutter: "stable",
          }}
          className="custom-scrollbar"
        >
          {previewQuestions.map((q, idx) => (
            <div key={idx} id={`question-${idx}`} style={{ marginBottom: 12 }}>
              <Card
                size="small"
                style={{
                  backgroundColor: selectedQuestions.has(idx)
                    ? "#E6F4FF"
                    : "#F5F5F5",
                  borderColor: selectedQuestions.has(idx)
                    ? "#1777FF"
                    : "#d9d9d9",
                  cursor: "pointer",
                  transition: "all 0.2s",
                }}
                onClick={() => handleSelectInternal(idx)}
              >
                <div style={{ display: "flex", alignItems: "flex-start" }}>
                  <Checkbox
                    checked={selectedQuestions.has(idx)}
                    style={{ marginRight: 16, flexShrink: 0, background: "" }}
                  />
                  <div style={{ flex: 1 }}>
                    <div
                      style={{
                        display: "flex",
                        alignItems: "center",
                        marginBottom: 12,
                      }}
                    >
                      <span
                        style={{
                          color: "#1777FF",
                          fontSize: 18,
                          fontWeight: 500,
                          marginRight: 8,
                        }}
                      >
                        题目{idx + 1}
                      </span>
                      <h4
                        style={{
                          fontSize: 18,
                          color: "#1a3353",
                          margin: 0,
                        }}
                      >
                        {q.title}
                      </h4>
                    </div>

                    <div style={{ marginLeft: 8 }}>
                      {q.answers.map((answer, i) => (
                        <div
                          key={i}
                          style={{
                            padding: "8px 0",
                            borderBottom: "1px solid #f0f0f0",
                            fontSize: 16,
                          }}
                        >
                          {answer}
                        </div>
                      ))}
                    </div>

                    <div style={{ marginTop: 12, textAlign: "right" }}>
                      <span
                        style={{
                          color: "#52c41a",
                          fontSize: 16,
                        }}
                      >
                        参考答案：{q.rights.join(", ")}
                      </span>
                    </div>
                  </div>
                </div>
              </Card>
            </div>
          ))}
        </div>

        {/* 底部滚动控制按钮 */}
        <div
          style={{
            padding: "12px 24px",
            backgroundColor: "rgba(240,248,255,0.8)",
            borderTop: "1px solid rgba(240,240,240,0.5)",
          }}
        >
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              gap: 16,
            }}
          >
            <Button
              shape="round"
              onClick={handleScrollToTop}
              style={{ flex: 1 }}
            >
              ↑ 回到顶部
            </Button>
            <Button
              shape="round"
              onClick={handleScrollToBottom}
              style={{ flex: 1 }}
            >
              ↓ 去到底部
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
};

export default QuestionCreateByAI;
