const axios = require('axios');

const baseUrl = 'https://search.censys.io/api/v2/hosts/search';
const headers = {
  'Accept': 'application/json',
  'Authorization': 'Basic YTcxYTg0NTAtZmU5NC00MjVjLWIwNTQtMzFkYjc3NWUyZDViOmFaVDJ6ejdoeFRVMVlqSFd6d2VnNzdDTHFKMGVjelVu'
};

function sendRequest(url) {
    axios.get(url, { headers })
        .then(response => {
        const { data } = response;
        const { result } = data;
        const { hits } = result;

        const ipList = hits.map(hit => hit.ip);
        sendIPsToGoServer(ipList);

        if (result.links && result.links.next) {
            const nextUrl = `${url}&cursor=${result.links.next}`;
            sendRequest(nextUrl);
        }
    })
    .catch(error => {
        console.error('Error:', error.message);
    });
}

function sendIPsToGoServer(ips) {
    const goServerUrl = 'http://localhost:8080/receive';
    axios.post(goServerUrl, { ips })
        .then(response => {
            console.log('IPs sent to Go server successfully');
        })
        .catch(error => {
            console.error('Error sending IPs to Go server:', error.message);
        });
}

const initialUrl = `${baseUrl}?q=%28FlareSolverr%29%20and%20services.service_name%3D%60HTTP%60&per_page=100&virtual_hosts=EXCLUDE`;
sendRequest(initialUrl);
