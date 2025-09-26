import React from 'react';
import { Routes, Route } from 'react-router-dom';
import { Container } from '@mui/material';
import Header from './components/Layout/Header';
import Footer from './components/Layout/Footer';
import HomePage from './pages/HomePage';
import MarketplacePage from './pages/MarketplacePage';
import NFTDetailPage from './pages/NFTDetailPage';
import UserProfilePage from './pages/UserProfilePage';
import CreateOrderPage from './pages/CreateOrderPage';
import CollectionsPage from './pages/CollectionsPage';
import ActivitiesPage from './pages/ActivitiesPage';
import TestPage from './pages/TestPage';
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
            <Route path="/collections" element={<CollectionsPage />} />
            <Route path="/activities" element={<ActivitiesPage />} />
            <Route path="/nft/:contract/:tokenId" element={<NFTDetailPage />} />
            <Route path="/profile/:address" element={<UserProfilePage />} />
            <Route path="/create-order" element={<CreateOrderPage />} />
            <Route path="/test" element={<TestPage />} />
          </Routes>
        </Container>
      </main>
      <Footer />
    </div>
  );
}

export default App;
