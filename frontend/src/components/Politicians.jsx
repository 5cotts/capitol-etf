import {useEffect, useState} from 'react';
import {
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableRow,
	Paper,
} from '@material-ui/core';


function politiciansTable(pidToNameMap) {
	return (
		<TableContainer component={Paper}>
		  <Table sx={{ minWidth: 650 }} aria-label="simple table">
			<TableBody>
			  {Object.entries(pidToNameMap).map(([pid, name]) => (
				<TableRow
				  key={pid}
				  sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
				>
				  <TableCell component="th" scope="row">
					<a href={`/trades/${pid}`}>{name}</a>
				  </TableCell>
				</TableRow>
			  ))}
			</TableBody>
		  </Table>
		</TableContainer>
	);
}

function Politicians() {
    const [error, setError] = useState(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [items, setItems] = useState([]);

    useEffect(() => {    
		fetch(`${process.env.REACT_APP_HILLTRADES_API_URL}/get-pids`)
			.then(res => res.json())
			.then(
				(result) => {
					setIsLoaded(true);
					setItems(result);
			},
			(error) => {
				setIsLoaded(true);
				setError(error);
			}
			)
    }, [])
  
    if (error) {
      return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
      return <div>Loading...</div>;
    } else {
	  return politiciansTable(items);
	}
  }

export default Politicians;
