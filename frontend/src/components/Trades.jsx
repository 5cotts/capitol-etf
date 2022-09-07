import {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import {
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableRow,
	Paper,
} from '@material-ui/core';


function tradesTable(trades) {
	return (
		<TableContainer component={Paper}>
		  <Table sx={{ minWidth: 650 }} aria-label="simple table">
			<TableBody>
			  {trades.map((trade) => (
				<TableRow
				  key={trade._txId}
				  sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
				>
				  <TableCell component="th" scope="row">
					{JSON.stringify(trade)}
				  </TableCell>
				</TableRow>
			  ))}
			</TableBody>
		  </Table>
		</TableContainer>
	);
}

function Trades() {
    const [error, setError] = useState(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [items, setItems] = useState([]);
    const params = useParams();

    useEffect(() => {
        fetch(`${process.env.REACT_APP_HILLTRADES_API_URL}/by-pid/${params.pid}`)
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
    }, [params.pid])

    if (error) {
        return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
        return <div>Loading...</div>;
    } else {
        return tradesTable(items);
    }
}

export default Trades;
