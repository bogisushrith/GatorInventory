import './App.css';
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { Root } from "./components/Root.tsx";
import { SignUp } from "./components/SignUp.tsx";
import { Login } from "./components/Login.tsx";
import { Dashboard } from "./components/Dashboard.tsx";
import { Add } from "./components/Add.tsx";
import { Update } from "./components/Update.tsx";
import { Users } from "./components/Users.tsx";
import { RequireAdmin } from "./components/RequireAdmin.tsx";

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Root/>}/>
                <Route path="login" element={<Login/>}/>
                <Route path="/signup" element={<SignUp/>}/>
                <Route path="/dashboard" element={<Dashboard/>}/>
                <Route path="/users" element={<RequireAdmin><Users/></RequireAdmin>}/>
                <Route path="/add-product" element={<RequireAdmin><Add/></RequireAdmin>}/>
                <Route path="/update-product" element={<RequireAdmin><Update/></RequireAdmin>}/>
            </Routes>
        </BrowserRouter>
    );
}

export default App;
