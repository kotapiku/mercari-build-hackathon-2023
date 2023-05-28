import React from "react";
import { Routes, Route, BrowserRouter } from "react-router-dom";
import { Home } from "./components/Home";
import { ItemDetail } from "./components/ItemDetail";
import { UserProfile } from "./components/UserProfile";
import { Listing } from "./components/Listing";
import "./App.css";
import { Header } from "./components/Header/Header";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Search } from "./components/Search";

export const App: React.VFC = () => {
  return (
    <>
      <ToastContainer position="bottom-center" />

      <BrowserRouter>
        <Header></Header>
        <Routes>
          <Route index element={<Home />} />
          <Route path="/item/:id" element={<ItemDetail />} />
          <Route path="/user/:id" element={<UserProfile />} />
          <Route path="/sell" element={<Listing />} />
          <Route path="/search" element={<Search />} />
        </Routes>
      </BrowserRouter>
    </>
  );
};
