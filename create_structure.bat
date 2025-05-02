@echo off
REM Create Journey Hub project structure

REM Create backend directory structure
mkdir backend\cmd\api
mkdir backend\internal\domain\user
mkdir backend\internal\domain\journey
mkdir backend\internal\domain\booking
mkdir backend\internal\domain\policy
mkdir backend\internal\domain\common

mkdir backend\internal\usecases\user
mkdir backend\internal\usecases\journey
mkdir backend\internal\usecases\booking
mkdir backend\internal\usecases\auth

mkdir backend\internal\interfaces\api\http\handlers
mkdir backend\internal\interfaces\api\http\dto
mkdir backend\internal\interfaces\api\grpc
mkdir backend\internal\interfaces\persistence\mongodb\models
mkdir backend\internal\interfaces\persistence\mongodb\repositories
mkdir backend\internal\interfaces\persistence\cache

mkdir backend\internal\infrastructure\middleware
mkdir backend\internal\infrastructure\server
mkdir backend\internal\infrastructure\database\mongodb
mkdir backend\internal\infrastructure\database\cache
mkdir backend\internal\infrastructure\config

mkdir backend\internal\pkg\logger
mkdir backend\internal\pkg\validator
mkdir backend\internal\pkg\metrics

mkdir backend\api\swagger
mkdir backend\config
mkdir backend\tests\unit\domain\user
mkdir backend\tests\unit\domain\journey
mkdir backend\tests\unit\domain\booking
mkdir backend\tests\unit\usecases\user
mkdir backend\tests\unit\usecases\journey
mkdir backend\tests\unit\usecases\booking
mkdir backend\tests\unit\interfaces\api
mkdir backend\tests\unit\interfaces\persistence
mkdir backend\tests\integration\api
mkdir backend\tests\integration\repository
mkdir backend\scripts\migrations

REM Create frontend directory structure
mkdir frontend\public
mkdir frontend\src\components\common\Button
mkdir frontend\src\components\common\Input
mkdir frontend\src\components\common\Card
mkdir frontend\src\components\layout\Header
mkdir frontend\src\components\layout\Footer

mkdir frontend\src\features\user\components
mkdir frontend\src\features\user\pages
mkdir frontend\src\features\user\services
mkdir frontend\src\features\journey\components
mkdir frontend\src\features\journey\pages
mkdir frontend\src\features\journey\services
mkdir frontend\src\features\booking\components
mkdir frontend\src\features\booking\pages
mkdir frontend\src\features\booking\services
mkdir frontend\src\features\auth\components
mkdir frontend\src\features\auth\pages
mkdir frontend\src\features\auth\services

mkdir frontend\src\hooks
mkdir frontend\src\services
mkdir frontend\src\utils

REM Create docs directory
mkdir docs

echo Project directory structure created successfully! 