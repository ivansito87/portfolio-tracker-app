import React, { useEffect, useState } from "react";
import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import { AppBar, Toolbar, Button, Container, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper } from "@mui/material";

interface Transaction {
  id: number;
  date: string;
  amount: number;
  type: string;
}

const Portfolio = () => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);

  useEffect(() => {
    fetch("http://localhost:8080/transactions")
      .then((response) => response.json())
      .then((data) => setTransactions(data))
      .catch((error) => console.error("Error fetching transactions:", error));
  }, []);

  return (
    <Container>
      <h1 className="text-2xl font-bold">Portfolio</h1>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Date</TableCell>
              <TableCell>Amount ($)</TableCell>
              <TableCell>Type</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {transactions.map((transaction, index) => (
              <TableRow key={index}>
                <TableCell>{transaction.date}</TableCell>
                <TableCell>{transaction.amount.toLocaleString()}</TableCell>
                <TableCell>{transaction.type}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
};

const Dashboard = () => <h1 className="text-2xl font-bold">Dashboard</h1>;
const Settings = () => <h1 className="text-2xl font-bold">Settings</h1>;

const Navbar = () => (
  <AppBar position="static">
    <Toolbar>
      <Button color="inherit" component={Link} to="/">Dashboard</Button>
      <Button color="inherit" component={Link} to="/portfolio">Portfolio</Button>
      <Button color="inherit" component={Link} to="/settings">Settings</Button>
    </Toolbar>
  </AppBar>
);

const App = () => {
  return (
    <Router>
      <Navbar />
      <Container>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/portfolio" element={<Portfolio />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Container>
    </Router>
  );
};

export default App;