import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import utc from 'dayjs/plugin/utc';
import weekday from 'dayjs/plugin/weekday';

dayjs.extend(relativeTime);
dayjs.extend(utc);
dayjs.extend(weekday);

interface DateProps {
  fromNow?: boolean;
  format?: string;
  special?: boolean;
  children?: string;
  stamp?: string;
}
export function Chrono({ fromNow, format = 'YYYY-MM-DD', special, stamp, children }: DateProps) {
  const raw = dayjs(stamp || children).utc();

  const process = () => {
    if (special !== undefined || special) {
      const today = dayjs().startOf('day').utc();
      const diff = raw.diff(today, 'days', true);
      // console.log('today: ', today.toString(), 'raw: ', raw.toString(), 'diff: ', diff);
      if (diff > -2 && diff <= -1) {
        return 'Yesterday';
      }
      if (diff > -1 && diff <= 0) {
        return 'Today';
      }
      if (diff > 0 && diff <= 1) {
        return 'Tomorrow';
      }
      if (diff > 1 && diff <= 7) {
        return dayjs()
          .weekday(today.weekday() + diff)
          .format('dddd');
      }
      return raw.format(format);
    }

    if (fromNow !== undefined || fromNow) {
      return raw.fromNow();
    }

    return raw.format(format);
  };

  return (
    <span className="dashotv-date" title={`${raw}`}>
      {process()}
    </span>
  );
}
