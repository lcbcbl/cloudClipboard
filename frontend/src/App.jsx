import React from 'react';
import { Divider, Card } from 'antd';
import { FileTextOutlined, CopyOutlined } from '@ant-design/icons';
import TextClipboard from './components/TextClipboard';
import FileManager from './components/FileManager';

function App() {
  return (
    <div className="app-container">
      <div className="app-header">
        <h1 className="app-title">云剪切板</h1>
        <p className="app-subtitle">轻松在设备间同步文本和文件</p>
      </div>
      
      <div className="app-content">
        {/* 字符串剪切板卡片 */}
        <div className="feature-column feature-column-text">
          <Card 
            className="feature-card"
            title={
              <div className="feature-title">
                <CopyOutlined className="feature-icon" />
                <span>字符串剪切板</span>
              </div>
            }
            bordered={false}
          >
            <TextClipboard />
          </Card>
        </div>
        
        {/* 文件管理卡片 */}
        <div className="feature-column feature-column-file">
          <Card 
            className="feature-card"
            title={
              <div className="feature-title">
                <FileTextOutlined className="feature-icon" />
                <span>文件剪切板</span>
              </div>
            }
            bordered={false}
          >
            <FileManager />
          </Card>
        </div>
      </div>
    </div>
  );
}

export default App;
