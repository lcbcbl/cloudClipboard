import React, { useState, useEffect } from 'react';
import { Button, Upload, List, Card, message, Space, Typography, Tooltip, Progress } from 'antd';
import { 
  DownloadOutlined, 
  DeleteOutlined, 
  UploadOutlined, 
  FileTextOutlined, 
  FileImageOutlined, 
  FilePdfOutlined, 
  FileWordOutlined, 
  FileExcelOutlined 
} from '@ant-design/icons';
import { fileApi } from '../services/api';

const { Title, Text } = Typography;

const FileManager = () => {
  const [fileList, setFileList] = useState([]);
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);

  // 加载所有文件
  const loadFiles = async () => {
    try {
      setLoading(true);
      const response = await fileApi.getAllFiles();
      // 确保response.data存在且files属性为数组，防止访问不存在的属性导致崩溃
      if (response && response.data) {
        setFileList(Array.isArray(response.data.files) ? response.data.files : []);
      } else {
        setFileList([]);
        message.error('服务器返回格式异常');
      }
    } catch (error) {
      // 提取后端返回的错误信息
      const errorMessage = error.response?.data?.message || '加载文件失败';
      message.error(errorMessage);
      console.error('Failed to load files:', error);
      setFileList([]); // 确保fileList始终为数组，防止页面崩溃
    } finally {
      setLoading(false);
    }
  };

  // 处理文件上传
  const handleFileUpload = async (file) => {
    setUploading(true);
    setUploadProgress(0);

    try {
      await fileApi.uploadFile(file);
      message.success('文件上传成功');
      loadFiles();
    } catch (error) {
      message.error('文件上传失败');
      console.error('Failed to upload file:', error);
    } finally {
      setUploading(false);
      setUploadProgress(0);
    }

    return false; // 阻止自动上传
  };

  // 处理文件下载
  const handleFileDownload = async (file) => {
    try {
      setLoading(true);
      await fileApi.downloadFile(file.id, file.filename);
      message.success('文件下载开始');
    } catch (error) {
      // 提取后端返回的错误信息
      const errorMessage = error.response?.data?.message || '文件下载失败';
      message.error(errorMessage);
      console.error('Failed to download file:', error);
    } finally {
      setLoading(false);
    }
  };

  // 处理文件删除
  const handleFileDelete = async (id) => {
    try {
      setLoading(true);
      await fileApi.deleteFile(id);
      message.success('文件删除成功');
      loadFiles();
    } catch (error) {
      // 提取后端返回的错误信息
      const errorMessage = error.response?.data?.message || '文件删除失败';
      message.error(errorMessage);
      console.error('Failed to delete file:', error);
    } finally {
      setLoading(false);
    }
  };

  // 格式化文件大小
  const formatFileSize = (size) => {
    if (size < 1024) {
      return `${size} B`;
    } else if (size < 1024 * 1024) {
      return `${(size / 1024).toFixed(2)} KB`;
    } else if (size < 1024 * 1024 * 1024) {
      return `${(size / (1024 * 1024)).toFixed(2)} MB`;
    } else {
      return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`;
    }
  };

  // 格式化时间
  const formatTime = (timestamp) => {
    return new Date(timestamp).toLocaleString('zh-CN');
  };

  // 判断文件类型并返回对应的图标或缩略图
  const getFileIcon = (file) => {
    const filename = file.filename.toLowerCase();
    const mimetype = file.mimetype;
    
    // 图片文件
    if (filename.endsWith('.jpg') || filename.endsWith('.jpeg') || filename.endsWith('.png') ||
        mimetype.startsWith('image/')) {
      // 这里应该使用文件的实际URL，暂时使用占位符
      return (
        <div style={{ 
          width: 40, 
          height: 40, 
          borderRadius: 4, 
          overflow: 'hidden', 
          marginRight: 8,
          display: 'inline-block',
          verticalAlign: 'middle',
          position: 'relative'
        }}>
          <img 
            src={`/api/files/${file.id}/thumbnail`} 
            alt={file.filename} 
            style={{ 
              width: '100%', 
              height: '100%', 
              objectFit: 'cover' 
            }} 
          />
          {/* 图片加载失败时，这个图标会显示 */}
          <div style={{ 
            position: 'absolute', 
            top: 0, 
            left: 0, 
            width: '100%', 
            height: '100%', 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: 'center', 
            backgroundColor: '#f5f5f5',
            fontSize: '24px',
            color: '#1890ff'
          }}>
            <FileImageOutlined />
          </div>
        </div>
      );
    }
    // PDF文件
    else if (filename.endsWith('.pdf') || mimetype === 'application/pdf') {
      return <FilePdfOutlined style={{ marginRight: 8, fontSize: 24, color: '#ff4d4f' }} />;
    }
    // Word文件
    else if (filename.endsWith('.docx') || filename.endsWith('.doc') || 
             mimetype === 'application/msword' || mimetype === 'application/vnd.openxmlformats-officedocument.wordprocessingml.document') {
      return <FileWordOutlined style={{ marginRight: 8, fontSize: 24, color: '#1890ff' }} />;
    }
    // Excel文件
    else if (filename.endsWith('.xlsx') || filename.endsWith('.xls') || 
             mimetype === 'application/vnd.ms-excel' || mimetype === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet') {
      return <FileExcelOutlined style={{ marginRight: 8, fontSize: 24, color: '#52c41a' }} />;
    }
    // 默认文件类型
    else {
      return <FileTextOutlined style={{ marginRight: 8, fontSize: 24, color: '#faad14' }} />;
    }
  };

  // 初始化加载
  useEffect(() => {
    loadFiles();
  }, []);

  return (
    <div>
      <div className="upload-section" style={{ width: '100%', textAlign: 'center' }}>
        <div style={{ marginTop: 0, width: '100%', display: 'flex', justifyContent: 'center' }}>
          <Upload
            beforeUpload={handleFileUpload}
            fileList={[]}
            showUploadList={false}
            accept="*/*"
            style={{ width: '85%' }}
          >
            <div className="upload-box" style={{ width: '100%' }}>
              <div className="upload-icon">
                <UploadOutlined style={{ fontSize: '32px', color: '#999' }} />
              </div>
              <div className="upload-text">
                <div>点击或拖拽文件到此处上传</div>
                <div style={{ fontSize: '12px', color: '#999', marginTop: 8 }}>
                  支持单个或多个文件上传
                </div>
              </div>
            </div>
          </Upload>
          {uploading && (
            <Progress
              percent={uploadProgress}
              status="active"
              style={{ marginTop: 16 }}
            />
          )}
        </div>
      </div>

      <div className="item-list">
        <Title level={4}>文件列表</Title>
        <List
          loading={loading}
          dataSource={fileList}
          locale={{ emptyText: '无数据' }}
          renderItem={(file) => (
            <Card className="item-card">
              <div className="item-header">
                <div style={{ display: 'flex', alignItems: 'center', flex: 1 }}>
                  <div style={{ display: 'flex', alignItems: 'center', flex: 1 }}>
                    {getFileIcon(file)}
                    <div style={{ flex: 1 }}>
                      <div className="item-title">
                        {file.filename}
                      </div>
                      <div className="item-meta">
                        <Space>
                          <Text>大小: {formatFileSize(file.size)}</Text>
                          <Text>上传时间: {formatTime(file.uploadTime)}</Text>
                          <Text>下载次数: {file.downloadCount}/{file.maxDownloads}</Text>
                        </Space>
                      </div>
                    </div>
                  </div>
                  <Space size="small">
                    <Button
                      icon={<DownloadOutlined />}
                      size="small"
                      type="text"
                      onClick={() => handleFileDownload(file)}
                      disabled={file.downloadCount >= file.maxDownloads}
                      title="下载"
                    />
                    <Button
                      icon={<DeleteOutlined />}
                      size="small"
                      type="text"
                      danger
                      onClick={() => handleFileDelete(file.id)}
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

export default FileManager;
