import axios from 'axios';

export function findTopTenTeaMatches(queryString) {
  return axios('http://127.0.0.1:5000/match', {
    method: 'POST',
    data: {
      userQuery: queryString,
    },
  });
}
