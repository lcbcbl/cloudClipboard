import React, { useState, useEffect } from 'react';
import { Button, Input, List, Card, message, Space, Typography, Tooltip } from 'antd';
import { CopyOutlined, DeleteOutlined, UploadOutlined } from '@ant-design/icons';
import { clipboardApi } from '../services/api';

const { TextArea } = Input;
const { Title, Text } = Typography;

const TextClipboard = () => {
  const [text, setText] = useState('');
  const [textItems, setTextItems] = useState([]);
  const [loading, setLoading] = useState(false);

  // 加载所有字符串
  const loadTextItems = async () => {
    try {
      setLoading(true);
      const response = await clipboardApi.getAllText();
      setTextItems(response.data.items);
    } catch (error) {
      const errorMessage = error.response?.data?.message || '加载字符串失败';
      message.error(errorMessage);
      console.error('Failed to load text items:', error);
    } finally {
      setLoading(false);
    }
  };

  // 上传字符串
  const handleUploadText = async () => {
    if (!text.trim()) {
      message.warning('请输入要上传的字符串');
      return;
    }

    try {
      setLoading(true);
      await clipboardApi.uploadText(text);
      message.success('字符串上传成功');
      setText('');
      loadTextItems();
    } catch (error) {
      const errorMessage = error.response?.data?.message || '字符串上传失败';
      message.error(errorMessage);
      console.error('Failed to upload text:', error);
    } finally {
      setLoading(false);
    }
  };

  // 删除字符串
  const handleDeleteText = async (id) => {
    try {
      setLoading(true);
      await clipboardApi.deleteTextById(id);
      message.success('字符串删除成功');
      loadTextItems();
    } catch (error) {
      const errorMessage = error.response?.data?.message || '字符串删除失败';
      message.error(errorMessage);
      console.error('Failed to delete text:', error);
    } finally {
      setLoading(false);
    }
  };

  // 复制到剪贴板
  const handleCopyText = (text) => {
    navigator.clipboard.writeText(text)
      .then(() => {
        message.success('已复制到剪贴板');
      })
      .catch(() => {
        message.error('复制失败');
      });
  };

  // 清空所有字符串
  const handleClearAllText = async () => {
    try {
      setLoading(true);
      await clipboardApi.clearAllText();
      message.success('已清空所有字符串');
      loadTextItems();
    } catch (error) {
      const errorMessage = error.response?.data?.message || '清空失败';
      message.error(errorMessage);
      console.error('Failed to clear all text:', error);
    } finally {
      setLoading(false);
    }
  };

  // 初始化加载
  useEffect(() => {
    loadTextItems();
  }, []);

  return (
    <div>
      <div className="upload-section">
        <TextArea
          rows={4}
          placeholder="请输入要粘贴的字符串，回车粘贴"
          value={text}
          onChange={(e) => setText(e.target.value)}
          onPressEnter={handleUploadText}
          loading={loading}
          autoSize={{ minRows: 3, maxRows: 6 }}
        />
      </div>

      <div className="item-list">
        {textItems.length > 0 && (
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
            <Button
              type="text"
              size="small"
              onClick={handleClearAllText}
              loading={loading}
              style={{ color: '#64748b' }}
            >
              全部清除
            </Button>
          </div>
        )}
        <List
            loading={loading}
            dataSource={textItems}
            locale={{ emptyText: '无数据' }}
            renderItem={(item) => (
              <Card className="item-card">
                <div className="item-content">
                  <div className="text-content">
                    <Text>{item.value}</Text>
                  </div>
                  <div className="item-actions">
                    <Space size="small">
                      <Button
                        icon={<CopyOutlined />}
                        size="small"
                        type="text"
                        onClick={() => handleCopyText(item.value)}
                        title="复制"
                      />
                      <Button
                        icon={<DeleteOutlined />}
                        size="small"
                        type="text"
                        danger
                        onClick={() => handleDeleteText(item.key)}
                        title="删除"
                      />
                    </Space>
                  </div>
                </div>
              </Card>
            )}
          />
      </div>
    </div>
  );
};

export default TextClipboard;
