import { sleep } from 'k6';
import k6perfana from 'k6/x/k6perfana';

export default function () {
  const startResponse = k6perfana.startPerfana();
  console.log(`startResponse ${startResponse.perfanaPayload}`);
  console.log(`startResponse ${startResponse.statusCode}`);
  console.log(`startResponse ${startResponse.body}`);


  sleep(66)

  const stopResponse = k6perfana.stopPerfana();
  console.log(`startResponse ${stopResponse.perfanaPayload}`);
  console.log(`startResponse ${stopResponse.statusCode}`);
  console.log(`startResponse ${stopResponse.body}`);

}