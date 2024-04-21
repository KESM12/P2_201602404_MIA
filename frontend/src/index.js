import * as React from "react";
import * as ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider, Route } from "react-router-dom";
import "./index.css";
import Nav from './pages/navbar';
import Card from './pages/card';
import Reports_Card from './pages/reports';
import No_Access from './pages/no_access';
import Modal from './pages/modal';
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap/dist/js/bootstrap.bundle.min.js'

function Home() {
  return (
    <div>
      <Nav />
      <div className='container'>
        <Card />
      </div>
      <Modal />
    </div>
  );
}

function Reports() {
  if(true){
    return (
      <div>
        <Nav />
        <div className='container'>
          <Reports_Card />
        </div>
        <Modal />
      </div>
    );
  } else {
    return (
      <div>
        <Nav />
        <div className='container'>
          <No_Access />
        </div>
        <Modal />
      </div>
    );
  }
}

const router = createBrowserRouter([
  {
    path: "/",
    element:<Home/>,
  },
  {
    path: "/reports",
    element:<Reports/>,
  },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);

export default router;