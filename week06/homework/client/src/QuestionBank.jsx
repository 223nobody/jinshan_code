import {
  Button,
  Radio,
  Input,
  Table,
  Dropdown,
  Menu,
  message,
  Pagination,
  Popconfirm,
  Space,
} from "antd";
import {
  PlusOutlined,
  DeleteOutlined,
  DownOutlined,
  ExclamationCircleOutlined,
} from "@ant-design/icons";
import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

const QuestionBank = () => {
  const [selectedType, setSelectedType] = useState("all");
  const [searchText, setSearchText] = useState("");
  const [selectedRowKeys, setSelectedRowKeys] = useState([]);
  const [statsData, setStatsData] = useState([]);
  const navigate = useNavigate();
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  // 根据题型获取数据的异步方法
  const fetchStats = async (type, pageParams) => {
    try {
      setLoading(true);
      let apiUrl = "/api/stats/summary";
      if (type === "single") apiUrl = "/api/stats/bytype1";
      if (type === "multiple") apiUrl = "/api/stats/bytype2";
      if (type === "programming") apiUrl = "/api/stats/bytype3";

      // 构建查询参数
      const params = new URLSearchParams({
        page: pageParams?.current || pagination.current,
        pageSize: pageParams?.pageSize || pagination.pageSize,
        search: pageParams?.search || searchText,
      });

      const response = await fetch(`${apiUrl}?${params}`);
      if (!response.ok) {
        throw new Error(`请求失败: ${response.status}`);
      }
      const result = await response.json();

      setPagination((prev) => ({
        ...prev,
        total: result.data.total || 0, // 更新总条数
      }));

      // 修改验证逻辑：允许空数组
      if (
        result.code !== 0 ||
        !result.data ||
        !Array.isArray(result.data.questions)
      ) {
        throw new Error("此时数据为空");
      }

      // 处理空数据情况
      const formattedData = result.data.questions.map((item) => ({
        key: item.id?.toString() || "", // 添加空值处理
        title: item.title || "无标题",
        type:
          item.type === 1
            ? "单选题"
            : item.type === 2
            ? "多选题"
            : item.type === 3
            ? "编程题"
            : "未知题型",
      }));

      setStatsData(formattedData);
      setError(null); // 清空错误状态
    } catch (err) {
      setStatsData([]); // 清空表格数据
      setError(err.message);
      message.error(err.message);
    } finally {
      setLoading(false);
    }
  };

  // 监听题型变化
  useEffect(() => {
    fetchStats(selectedType, {
      current: pagination.current,
      pageSize: pagination.pageSize,
    });
  }, [selectedType, pagination.current, pagination.pageSize]); // 添加分页参数依赖

  //删除操作处理函数
  const handleDelete = async (ids) => {
    try {
      setLoading(true);
      const response = await fetch("/api/questions/batch-delete", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ ids: ids.map(Number) }),
      });

      const result = await response.json();
      // 删除后重新获取数据
      fetchStats(selectedType, {
        current: pagination.current,
        pageSize: pagination.pageSize,
      });

      // 同步更新总条数
      setPagination((prev) => ({
        ...prev,
        total: prev.total - ids.length, // 立即减少总数
      }));

      if (!response.ok || result.code !== 0) {
        throw new Error(result.msg || `删除失败: ${response.status}`);
      }
      setStatsData((prev) => prev.filter((item) => !ids.includes(item.key)));
      setSelectedRowKeys([]);
      message.success(`成功删除${ids.length}条记录`);
    } catch (err) {
      message.error(err.message);
    } finally {
      setLoading(false);
    }
  };

  // 处理菜单点击
  const handleMenuClick = (e) => {
    switch (e.key) {
      case "ai":
        navigate("/ai-create"); // 跳转到AI出题页面
        break;
      case "manual":
        navigate("/manual-create"); // 跳转到自主出题页面
        break;
      default:
        break;
    }
  };
  // 出题按钮菜单配置
  const createMenu = (
    <Menu onClick={handleMenuClick}>
      <Menu.Item key="ai">AI出题</Menu.Item>
      <Menu.Item key="manual">自主出题</Menu.Item>
    </Menu>
  );

  // 表格列配置
  const columns = [
    {
      title: "题目",
      dataIndex: "title",
      width: 800,
      // 允许内容换行显示
      onCell: () => ({
        style: { whiteSpace: "normal", wordBreak: "break-word" },
      }), // 超出换行
      render: (text) => <span style={{ color: "	#1E90FF" }}>{text}</span>,
    },
    {
      title: "题型",
      align: "center",
      width: 200,
      dataIndex: "type",
      render: (type) => <span style={{ color: "black" }}>{type}</span>,
    },
    {
      title: "操作",
      align: "center",
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            style={{ color: "	#1E90FF" }}
            onClick={() => navigate(`/manual-create/${record.key}`)}
          >
            编辑
          </Button>
          {/* 单条删除：包裹 Popconfirm */}
          <Popconfirm
            title={<span style={{ fontWeight: 500 }}>确认要删除此题？</span>}
            icon={<ExclamationCircleOutlined style={{ color: "#faad14" }} />}
            okText="删除"
            cancelText="取消"
            okButtonProps={{
              style: {
                width: 80,
                background: "#d9001b",
                borderColor: "#d9001b",
                borderRadius: 4,
              },
            }}
            cancelButtonProps={{ style: { width: 80 } }}
            onConfirm={() => handleDelete([record.key])}
            placement="topRight"
            overlayStyle={{ borderRadius: 8 }}
          >
            <Button type="link" danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      {/* 增加加载状态指示 */}
      {loading && <div style={{ margin: 20 }}>数据加载中...</div>}

      {/* 显示错误信息 */}
      {error && <div style={{ margin: 20, color: "red" }}>错误: {error}</div>}

      <div
        style={{
          display: "flex",
          alignItems: "center",
          gap: 16,
          flexWrap: "nowrap",
          marginBottom: 16,
        }}
      >
        <div
          style={{
            display: "flex",
            alignItems: "center",
            flexShrink: 0,
          }}
        >
          <span style={{ marginRight: 12 }}>题型：</span>
          <Radio.Group
            value={selectedType}
            onChange={(e) => {
              setSelectedType(e.target.value);
              setPagination({
                current: 1,
                pageSize: 10,
                total: pagination.total,
              });
            }}
          >
            <Radio.Button value="all">全部</Radio.Button>
            <Radio.Button value="single">单选题</Radio.Button>
            <Radio.Button value="multiple">多选题</Radio.Button>
            <Radio.Button value="programming">编程题</Radio.Button>
          </Radio.Group>
        </div>

        <Input.Search
          placeholder="请输入试题名称"
          style={{
            flex: 1,
            minWidth: 200,
          }}
          allowClear
          onSearch={(value) => {
            setSearchText(value);
            setPagination((prev) => ({
              ...prev,
              current: 1,
            }));
            fetchStats(selectedType, {
              current: 1,
              pageSize: pagination.pageSize,
              search: value,
            });
          }}
          onChange={(e) => setSearchText(e.target.value)}
        />
      </div>

      {/* 操作按钮 */}
      <div
        style={{
          display: "flex",
          justifyContent: "right",
          marginBottom: 16,
        }}
      >
        <div style={{ marginRight: 8 }}>
          <Dropdown overlay={createMenu} placement="bottomRight">
            <Button
              style={{
                background: "#1777FF",
                borderRadius: "4px",
                height: 32,
                padding: "4px 15px",
                display: "inline-flex",
                alignItems: "center",
                boxShadow: "0 2px 0 rgba(5, 145, 255, 0.1)",
              }}
            >
              <PlusOutlined style={{ fontSize: 14, color: "#fff" }} />
              <span style={{ margin: "0 8px", color: "#fff" }}>出题</span>
              <DownOutlined
                style={{
                  fontSize: 12,
                  color: "#fff",
                  transform: "translateY(1px)",
                }}
              />
            </Button>
          </Dropdown>
        </div>
        {/* 批量删除：包裹 Popconfirm */}
        <Popconfirm
          title={
            <span style={{ fontWeight: 500 }}>
              确认删除 {selectedRowKeys.length} 项？
            </span>
          }
          icon={<ExclamationCircleOutlined style={{ color: "#faad14" }} />}
          okText="删除"
          cancelText="取消"
          okButtonProps={{
            style: {
              width: 80,
              background: "#d9001b",
              borderColor: "#d9001b",
            },
          }}
          cancelButtonProps={{ style: { width: 80 } }}
          onConfirm={() => handleDelete(selectedRowKeys)}
          disabled={!selectedRowKeys.length}
        >
          <Button
            danger
            icon={<DeleteOutlined />}
            disabled={!selectedRowKeys.length}
            style={{
              marginLeft: 8,
              background: "#d9001b",
              borderColor: "#d9001b",
              color: "#fff",
            }}
          >
            批量删除
          </Button>
        </Popconfirm>
      </div>

      {/* 试题表格 */}
      <Table
        columns={columns}
        dataSource={statsData}
        pagination={false}
        rowSelection={{
          selectedRowKeys,
          onChange: (keys) => setSelectedRowKeys(keys),
        }}
        bordered
        style={{
          border: "1px solid #f0f0f0",
          borderRadius: 4,
        }}
      />
      {/* 分页组件 */}
      <div
        style={{
          display: "flex",
          justifyContent: "flex-end",
        }}
      >
        <Pagination
          current={pagination.current}
          pageSize={pagination.pageSize}
          total={pagination.total}
          showSizeChanger
          showQuickJumper
          pageSizeOptions={["10", "20", "50", "100"]}
          showTotal={(total) => `共 ${total} 条`}
          onChange={(page, pageSize) => {
            setPagination((prev) => ({
              ...prev,
              current: page,
              pageSize,
            }));
          }}
          onShowSizeChange={(current, size) => {
            setPagination((prev) => ({
              ...prev,
              current: 1, // 切换分页大小时重置到第一页
              pageSize: size,
            }));
          }}
        />
      </div>
    </div>
  );
};

export default QuestionBank;
