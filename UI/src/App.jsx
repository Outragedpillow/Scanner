import logo from './logo.svg';
import styles from './App.module.css';
import HomePage from "./pages/Home.jsx"

function App() {
  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <HomePage />
      </header>
    </div>
  );
}

export default App;
