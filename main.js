import { sleep } from 'k6';
import k6perfana from 'k6/x/k6perfana';

export default function () {
  k6perfana.startPerfana();

  sleep(66)

  k6perfana.stopPerfana();
}