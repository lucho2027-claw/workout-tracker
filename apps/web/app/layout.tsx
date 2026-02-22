export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body style={{ fontFamily: 'sans-serif', maxWidth: 760, margin: '40px auto', padding: '0 16px' }}>
        <h1>Workout Tracker</h1>
        {children}
      </body>
    </html>
  );
}
