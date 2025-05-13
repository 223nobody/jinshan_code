import React, { useState, useEffect } from "react";
import { Layout, Menu, Card, ConfigProvider, Spin } from "antd";
import {
  FormOutlined,
  DatabaseOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined,
} from "@ant-design/icons";
import { BrowserRouter as Router, Route, Link, Routes } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import QuestionBank from "./QuestionBank";
import QuestionCreateByHand from "./QuestionCreateByHand";
import QuestionCreateByAI from "./QuestionCreateByAI";
import QuestionEdit from "./QuestionEdit";
import readmecontent from "../../readme.md";
import image from "./images/image.png";

const { Header, Sider, Content } = Layout;

const StudyNotes = () => {
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(readmecontent)
      .then((response) => response.text())
      .then((text) => {
        setContent(text);
        setLoading(false);
      })
      .catch((error) => {
        console.error("读取文件失败:", error);
        setLoading(false);
      });
  }, []);

  return (
    <Card
      bordered
      style={{
        border: "1px solid #1777FF",
        minHeight: 300,
        margin: 24,
        borderRadius: 8,
      }}
    >
      {loading ? (
        <div style={{ textAlign: "center", padding: 40 }}>
          <Spin tip="正在加载学习心得..." />
        </div>
      ) : (
        <div
          style={{
            padding: 24,
            fontFamily: "'Segoe UI', 'PingFang SC', sans-serif",
            lineHeight: 1.6,
          }}
        >
          <ReactMarkdown
            components={{
              h1: ({ node, ...props }) => (
                <h1 style={{ color: "#1777FF" }} {...props}>
                  {props.children || "Default Heading"}
                </h1>
              ),
              code: ({ node, inline, className, children, ...props }) => (
                <code
                  style={{
                    background: "#f6f8fa",
                    padding: "2px 4px",
                    borderRadius: 4,
                    fontFamily: "Menlo, Monaco, Consolas",
                  }}
                  {...props}
                >
                  {children}
                </code>
              ),
            }}
          >
            {content}
          </ReactMarkdown>
        </div>
      )}
    </Card>
  );
};

const App = () => {
  const [collapsed, setCollapsed] = useState(false);

  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#1777FF",
          fontSize: 16,
          lineHeight: 1.6,
        },
        components: {
          Menu: {
            itemSelectedBg: "#1777FF",
            itemSelectedColor: "#fff",
            itemHoverBg: "rgba(23, 119, 255, 0.1)",
            itemBorderRadius: 4,
          },
        },
      }}
    >
      <Router>
        <Layout style={{ minHeight: "100vh" }}>
          <Header
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
              background: "#001629",
              padding: "0 24px",
              boxShadow: "0 2px 8px #001629",
            }}
          >
            <div
              style={{
                color: "white",
                fontWeight: "bold",
                fontSize: 30,
              }}
            >
              <img src={image} alt="Logo" />
            </div>
            <div style={{ color: "white", fontSize: 20, marginRight: "auto" }}>
              武汉科技大学 付坤 大作业
            </div>
          </Header>

          <Layout>
            <Sider
              width={200}
              collapsedWidth={80}
              collapsible
              collapsed={collapsed}
              onCollapse={(value) => setCollapsed(value)}
              style={{
                background: "#001629",
                boxShadow: "2px 0 8px #001629",
              }}
            >
              <Menu
                mode="inline"
                defaultSelectedKeys={["1"]}
                inlineCollapsed={collapsed}
                style={{
                  borderRight: 0,
                  background: "#001629",
                  padding: "8px 0",
                }}
              >
                <Menu.Item
                  key="1"
                  icon={<FormOutlined />}
                  style={{ color: "white", margin: "8px 0" }}
                >
                  <Link to="/">学习心得</Link>
                </Menu.Item>
                <Menu.Item
                  key="2"
                  icon={<DatabaseOutlined />}
                  style={{ color: "white", margin: "8px 0" }}
                >
                  <Link to="/questions">题库管理</Link>
                </Menu.Item>

                <Menu.Item
                  key="3"
                  onClick={() => setCollapsed(!collapsed)}
                  style={{ color: "white", margin: "8px 0", height: 0 }}
                >
                  <div style={{ color: "white" }}>
                    {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                  </div>
                </Menu.Item>
              </Menu>
            </Sider>

            <Content style={{ overflow: "auto" }}>
              <Routes>
                <Route path="/" element={<StudyNotes />} />
                <Route path="/questions" element={<QuestionBank />} />
                <Route
                  path="/manual-create"
                  element={<QuestionCreateByHand />}
                />
                <Route path="/manual-create/:id" element={<QuestionEdit />} />
                <Route path="/ai-create" element={<QuestionCreateByAI />} />
              </Routes>
            </Content>
          </Layout>
        </Layout>
      </Router>
    </ConfigProvider>
  );
};

export default App;
