import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  timeout: 30000
});

// 字符串剪切板API
export const clipboardApi = {
  // 上传字符串
  uploadText: (text) => {
    return api.post('/clipboard/text', { text });
  },
  
  // 获取所有字符串
  getAllText: () => {
    return api.get('/clipboard/text');
  },
  
  // 获取指定字符串
  getTextById: (id) => {
    return api.get(`/clipboard/text/${id}`);
  },
  
  // 删除指定字符串
  deleteTextById: (id) => {
    return api.delete(`/clipboard/text/${id}`);
  },
  
  // 清空所有字符串
  clearAllText: () => {
    return api.delete('/clipboard/text');
  }
};

// 文件API
export const fileApi = {
  // 上传文件
  uploadFile: (file) => {
    const formData = new FormData();
    formData.append('file', file);
    return api.post('/files', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
        console.log(`Upload progress: ${percentCompleted}%`);
      }
    });
  },
  
  // 获取所有文件
  getAllFiles: () => {
    return api.get('/files');
  },
  
  // 获取文件信息
  getFileInfo: (id) => {
    return api.get(`/files/${id}`);
  },
  
  // 下载文件
  downloadFile: (id, filename) => {
    return api.get(`/files/${id}/download`, {
      responseType: 'blob',
      onDownloadProgress: (progressEvent) => {
        if (progressEvent.total) {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          console.log(`Download progress: ${percentCompleted}%`);
        }
      }
    }).then((response) => {
      // 创建下载链接
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', filename);
      document.body.appendChild(link);
      link.click();
      
      // 清理
      link.parentNode.removeChild(link);
      window.URL.revokeObjectURL(url);
    });
  },
  
  // 删除文件
  deleteFile: (id) => {
    return api.delete(`/files/${id}`);
  }
};

// 健康检查API
export const healthApi = {
  checkHealth: () => {
    return api.get('/health');
  }
};
