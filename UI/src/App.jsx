import { createSignal, onMount, onCleanupP } from 'solid-js';

function App() {
  // Define a signal to store the data
  const [data, setData] = createSignal([]);

  // Function to fetch data from the API
  const fetchData = async () => {
    try {
      const response = await fetch('http://127.0.0.1:1234/api/data');
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const jsonData = await response.json();
      setData(jsonData.currentsignouts); // Use the correct property name here
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };

  // Fetch data initially
  onMount(fetchData);

  // Set up an interval to fetch data every 2 seconds
  onCleanupP(() => {
    const intervalId = setInterval(fetchData, 2000);
    return () => clearInterval(intervalId);
  });

  return (
    <div>
      <h1>Current Sign Outs</h1>
      <ul>
        {data().map((item, index) => (
          <li key={index}>{item}</li>
        ))}
      </ul>    
    </div>
  );
}

export default App;

