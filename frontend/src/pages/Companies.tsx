import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import api from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';

interface Company {
  ID: number;
  name: string;
  address: string;
}

const companySchema = z.object({
  name: z.string().min(1, "Name is required"),
  address: z.string().optional(),
});

type CompanyFormValues = z.infer<typeof companySchema>;

export default function Companies() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [isCreating, setIsCreating] = useState(false);

  const form = useForm<CompanyFormValues>({
    resolver: zodResolver(companySchema),
    defaultValues: {
      name: '',
      address: '',
    },
  });

  const fetchCompanies = async () => {
      try {
        const res = await api.get('/api/companies');
        setCompanies(res.data);
      } catch (error) {
        console.error("Failed to fetch companies", error);
      }
    };

  useEffect(() => {
    fetchCompanies();
  }, []);

  const onSubmit = async (data: CompanyFormValues) => {
      try {
          await api.post('/api/companies', data);
          await fetchCompanies();
          setIsCreating(false);
          form.reset();
      } catch (error) {
          console.error("Failed to create company", error);
      }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold tracking-tight">Companies</h2>
        <Button onClick={() => setIsCreating(!isCreating)}>{isCreating ? 'Cancel' : 'Add Company'}</Button>
      </div>
      
      {isCreating && (
          <Card>
              <CardHeader><CardTitle>New Company</CardTitle></CardHeader>
              <CardContent>
                  <Form {...form}>
                      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                          <FormField control={form.control} name="name" render={({ field }) => (
                              <FormItem>
                                  <FormLabel>Name</FormLabel>
                                  <FormControl><Input {...field} /></FormControl>
                                  <FormMessage />
                              </FormItem>
                          )} />
                          <FormField control={form.control} name="address" render={({ field }) => (
                              <FormItem>
                                  <FormLabel>Address</FormLabel>
                                  <FormControl><Input {...field} /></FormControl>
                                  <FormMessage />
                              </FormItem>
                          )} />
                          <Button type="submit">Save</Button>
                      </form>
                  </Form>
              </CardContent>
          </Card>
      )}

      <Card>
        <CardHeader>
          <CardTitle>All Companies</CardTitle>
        </CardHeader>
        <CardContent>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>ID</TableHead>
                        <TableHead>Name</TableHead>
                        <TableHead>Address</TableHead>
                        <TableHead>Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {companies.map((company) => (
                        <TableRow key={company.ID}>
                            <TableCell>{company.ID}</TableCell>
                            <TableCell>{company.name}</TableCell>
                            <TableCell>{company.address}</TableCell>
                            <TableCell>
                                <Button variant="ghost" size="sm">Edit</Button>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </CardContent>
      </Card>
    </div>
  );
}
