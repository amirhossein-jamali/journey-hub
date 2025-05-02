# راهنمای رمزنگاری تنظیمات

این مستند نحوه استفاده از سیستم رمزنگاری تنظیمات در محیط‌های مختلف پروژه Journey Hub را توضیح می‌دهد.

## مقدمه

سیستم رمزنگاری تنظیمات برای محافظت از داده‌های حساس مانند رمزهای عبور و کلیدهای دسترسی در فایل‌های پیکربندی طراحی شده است.

## متغیرهای محیطی کلیدی

سیستم رمزنگاری از متغیرهای محیطی زیر استفاده می‌کند:

- `JH_CONFIG_ENCRYPTION_KEY`: کلید رمزنگاری برای رمزگذاری و رمزگشایی مقادیر
- `JH_SKIP_CONFIG_ENCRYPTION`: اگر `true` باشد، سیستم رمزنگاری غیرفعال می‌شود (برای محیط‌های CI/CD)
- `JH_ENV`: محیط اجرای برنامه (`development`, `test`, `staging`, `production`)

## محیط‌های مختلف

### محیط توسعه (Development)

در محیط توسعه، اگر `JH_CONFIG_ENCRYPTION_KEY` تنظیم نشده باشد، سیستم از یک کلید پیش‌فرض استفاده می‌کند. این حالت فقط برای توسعه‌دهندگان مناسب است.

برای ایجاد مقادیر رمزنگاری شده:
```bash
cd backend
./scripts/setup-encrypted-configs.sh
```

### محیط تست (Test)

برای تست‌های محلی، می‌توانید `JH_CONFIG_ENCRYPTION_KEY` را تنظیم کنید یا از کلید پیش‌فرض استفاده کنید.

برای تست‌های CI/CD، تنظیم `JH_SKIP_CONFIG_ENCRYPTION=true` توصیه می‌شود.

### محیط تولید (Production)

در محیط تولید، **باید** `JH_CONFIG_ENCRYPTION_KEY` را تنظیم کنید. استفاده از کلید پیش‌فرض در محیط تولید خطای اجرا ایجاد می‌کند.

```bash
export JH_ENV=production
export JH_CONFIG_ENCRYPTION_KEY="your-secure-key-here"
```

## محیط CI/CD

در محیط‌های CI/CD، توصیه می‌شود سیستم رمزنگاری را با تنظیم متغیر محیطی `JH_SKIP_CONFIG_ENCRYPTION=true` غیرفعال کنید:

```yaml
env:
  JH_SKIP_CONFIG_ENCRYPTION: "true"
```

همچنین می‌توانید فایل‌های پیکربندی مخصوص CI/CD بدون مقادیر رمزنگاری شده ایجاد کنید.

## دستورات مفید

### رمزنگاری یک مقدار جدید

```bash
cd backend
go run cmd/encrypt/main.go -value "your-secret-value"
```

### راه‌اندازی خودکار مقادیر رمزنگاری شده

```bash
cd backend
./scripts/setup-encrypted-configs.sh
```

### بررسی عملکرد hot-reload تنظیمات

```bash
cd backend
./scripts/test-config-hot-reload.sh
```

## توجه امنیتی

1. هرگز کلیدهای رمزنگاری را در مخزن گیت ذخیره نکنید
2. هرگز اسکریپت `setup-encrypted-configs.sh` را در محیط CI/CD اجرا نکنید
3. برای محیط تولید، از یک سیستم مدیریت راز (secret management) مثل Vault یا AWS KMS استفاده کنید 