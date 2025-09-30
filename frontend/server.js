import express from 'express';
import { createProxyMiddleware } from 'http-proxy-middleware';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const app = express();


// API代理到Go后端
app.use('/api', (req, res, next) => {
  console.log('Proxying:', req.method, req.url);
  next();
}, createProxyMiddleware({
  target: 'http://localhost:8080/api',
  changeOrigin: true,
}));
// 静态资源服务
app.use(express.static(path.join(__dirname, 'dist')));

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Express server running at http://localhost:${PORT}`);
});
