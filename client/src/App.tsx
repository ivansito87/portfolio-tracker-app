import React, { useEffect, useState } from "react";
import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import { AppBar, Toolbar, Button, Container, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, CircularProgress, Alert } from "@mui/material";

interface Transaction {
  date: string;
  amount: number;
  type: string;
}

const Portfolio = () => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("http://3.145.169.224:8080/api/transactions")
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        setTransactions(data);
        setLoading(false);
      })
      .catch((error) => {
        setError(error.message);
        setLoading(false);
      });
  }, []);

  return (
    <Container>
      <h1 className="text-2xl font-bold">Portfolio</h1>
      {loading && <CircularProgress />}
      {error && <Alert severity="error">{error}</Alert>}
      {!loading && !error && (
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
      )}
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