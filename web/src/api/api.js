import axios from 'axios'

export default {
    get: function(url, params = {}, alert_error = true) {
        return new Promise((resolve, reject) => {
            axios.get(url, {
                url: params
            }).then(response => {
                resolve(response.data);
            }).catch(error => {
                if (alert_error) {
                    alert('GET: ' + url + '\n' + JSON.stringify(error));
                }
                reject(error);
            });
        });
    },
    post: function(url, params = {}) {
        return new Promise((resolve, reject) => {
            axios.post(url, params).then(response => {
                resolve(response.data);
            }).catch(error => {
                alert('POST: ' + url + '\n' + JSON.stringify(error));
                reject(error);
            });
        });
    },
    put: function(url, params = {}) {
        return new Promise((resolve, reject) => {
            axios.put(url, params).then(response => {
                resolve(response.data);
            }).catch(error => {
                alert('PUT: ' + url + '\n' + JSON.stringify(error));
                reject(error);
            });
        });
    }
}
