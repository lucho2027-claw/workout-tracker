import { Card } from "../components/ui/card";
import { theme } from "../lib/theme";

export default function Home() {
  return (
    <Card>
      <h2 className={theme.title}>Welcome</h2>
      <p className={theme.subtitle + " mt-2"}>Track exercises, sets, reps, and weight with a clean workflow.</p>
    </Card>
  );
}
