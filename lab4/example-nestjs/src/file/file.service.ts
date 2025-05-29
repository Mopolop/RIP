import { Injectable } from '@nestjs/common';
import * as fs from 'fs';
import * as path from 'path';

@Injectable()
export class FileService<T> {
    private readonly filePath: string;

    constructor(filename: string) {
        this.filePath = path.join(process.cwd(), 'src/assets', filename);
        this.ensureFileExists();
    }

    private ensureFileExists(): void {
        if (!fs.existsSync(this.filePath)) {
            fs.writeFileSync(this.filePath, '[]', 'utf8');
        }
    }

    read(): T {
        const data = fs.readFileSync(this.filePath, 'utf8');
        return JSON.parse(data);
    }

    write(data: T): void {
        fs.writeFileSync(this.filePath, JSON.stringify(data, null, 2), 'utf8');
    }

    add(item: any): void {
        const data = this.read() as any[];
        data.push(item);
        this.write(data as T);
    }
}