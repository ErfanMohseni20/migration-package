# migration-package

یک بسته ساده برای اجرای مهاجرت‌های پایگاه داده PostgreSQL با استفاده از `golang-migrate`.

## امکانات

- اجرای مهاجرت‌ها با `up`
- بازگرداندن مهاجرت‌ها با `down`
- تنظیم نسخه با `force`
- نمایش وضعیت فعلی با `status`

## پیش‌نیازها

- Go 1.24+
- PostgreSQL

## پیکربندی

فایل `.env` یا متغیرهای محیطی باید شامل یکی از موارد زیر باشد:

- `MIGRATION_URL` برای اتصال مستقیم به دیتابیس
- یا `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_DATABASE`, `DB_SSLMODE`

همچنین مسیر فایل‌های مهاجرت را می‌توان با `MIGRATIONS_PATH` مشخص کرد.

## اجرا

```bash
./cmd/migrate/migrate -action up
./cmd/migrate/migrate -action down -steps 1
./cmd/migrate/migrate -action force -version 2
./cmd/migrate/migrate -action status
```

اگر از باینری ساخته شده استفاده می‌کنید، نام آن را در دستورها جایگزین کنید.

## مسیر پیش‌فرض مهاجرت

- `internal/db/migrations`

## توجه

- `force` نیاز دارد `-version` تنظیم شده باشد.
- اگر `MIGRATION_URL` تنظیم نشده باشد، پروژه از متغیرهای `DB_*` برای ساخت URL استفاده می‌کند.
