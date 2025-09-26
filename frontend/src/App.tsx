import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { Container } from '@mui/material';
import Header from './components/Layout/Header.tsx';
import Footer from './components/Layout/Footer.tsx';
import HomePage from './pages/HomePage';
import MarketplacePage from './pages/MarketplacePage';
import NFTDetailPage from './pages/NFTDetailPage';
import UserProfilePage from './pages/UserProfilePage';
import CreateOrderPage from './pages/CreateOrderPage';
import './App.css';

function App() {
  return (
    <div className="App">
      <Header />
      <main style={{ minHeight: 'calc(100vh - 64px - 60px)', paddingTop: '20px', paddingBottom: '20px' }}>
        <Container maxWidth="xl">
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/marketplace" element={<MarketplacePage />} />
            <Route path="/nft/:contract/:tokenId" element={<NFTDetailPage />} />
            <Route path="/profile/:address" element={<UserProfilePage />} />
            <Route path="/create-order" element={<CreateOrderPage />} />
          </Routes>
        </Container>
      </main>
      <Footer />
    </div>
  );
}

export default App;
