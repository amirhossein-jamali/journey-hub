import React from 'react';
import { Link } from 'react-router-dom';
import './HomePage.css';

const HomePage = () => {
  return (
    <div className="home-container">
      <div className="hero-section">
        <h1>خوش آمدید به Journey Hub</h1>
        <p>مقصد بعدی سفر شما را با ما کشف کنید</p>
        <div className="cta-buttons">
          <Link to="/journeys" className="btn primary-btn">جستجوی سفرها</Link>
          <Link to="/register" className="btn secondary-btn">ثبت نام</Link>
        </div>
      </div>

      <div className="features-section">
        <div className="feature-card">
          <i className="fas fa-ticket-alt"></i>
          <h3>رزرو آسان</h3>
          <p>به سادگی بلیط سفر خود را رزرو کنید</p>
        </div>
        <div className="feature-card">
          <i className="fas fa-map-marked-alt"></i>
          <h3>مقصدهای متنوع</h3>
          <p>صدها مقصد جذاب در انتظار شماست</p>
        </div>
        <div className="feature-card">
          <i className="fas fa-credit-card"></i>
          <h3>پرداخت امن</h3>
          <p>پرداخت سریع و امن با انواع روش‌های پرداخت</p>
        </div>
      </div>

      <div className="popular-journeys">
        <h2>سفرهای محبوب</h2>
        <div className="journey-cards">
          {/* نمونه کارت‌های سفر - در نسخه واقعی از API دریافت می‌شوند */}
          <div className="journey-card">
            <div className="journey-image"></div>
            <h3>تهران به مشهد</h3>
            <p>شروع از <span className="price">۲,۵۰۰,۰۰۰</span> تومان</p>
            <Link to="/journey/1" className="btn small-btn">مشاهده جزئیات</Link>
          </div>
          <div className="journey-card">
            <div className="journey-image"></div>
            <h3>تهران به اصفهان</h3>
            <p>شروع از <span className="price">۱,۸۰۰,۰۰۰</span> تومان</p>
            <Link to="/journey/2" className="btn small-btn">مشاهده جزئیات</Link>
          </div>
          <div className="journey-card">
            <div className="journey-image"></div>
            <h3>تهران به شیراز</h3>
            <p>شروع از <span className="price">۲,۲۰۰,۰۰۰</span> تومان</p>
            <Link to="/journey/3" className="btn small-btn">مشاهده جزئیات</Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default HomePage; 