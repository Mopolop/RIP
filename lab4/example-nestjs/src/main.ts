import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import * as cors from 'cors';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  // Включение CORS
  app.use(cors());

  await app.listen(3000);
  console.log('Сервер запущен на http://localhost:3000');
}
bootstrap();
