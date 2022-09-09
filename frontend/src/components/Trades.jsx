import {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import {
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableRow,
  TableHead,
	Paper,
} from '@material-ui/core';


function tradesTable(trades) {
	return (
		<TableContainer component={Paper}>
		  <Table sx={{ minWidth: 650 }} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>Politician</TableCell>
            <TableCell>Traded Issuer</TableCell>
            <TableCell>Filing Date</TableCell>
            <TableCell>Trade Date</TableCell>
            <TableCell>Owner</TableCell>
            <TableCell>Transaction Type</TableCell>
            <TableCell>Asset Type</TableCell>
            <TableCell>Size</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {trades.map((trade) => (
            <TableRow
              key={trade._txId}
              sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
              <TableCell component="th" scope="row">{`${trade.politician.firstName} ${trade.politician.lastName}`}</TableCell>
              <TableCell>
                <a href={trade.filingURL}>{`${trade.issuer.issuerName} (${trade.issuer.issuerTicker})`}</a>
              </TableCell>
              <TableCell>{trade.filingDate}</TableCell>
              <TableCell>{trade.txDate}</TableCell>
              <TableCell>{trade.owner}</TableCell>
              <TableCell>{trade.txType}</TableCell>
              <TableCell>
                {trade.asset.assetType === "stock-options" ? trade.asset.instrument : "Stock"}
              </TableCell>
              <TableCell>
                ${trade.asset.assetType === "stock-options" ?  trade.value : `${Math.floor(trade.price * trade.size)}`}
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
