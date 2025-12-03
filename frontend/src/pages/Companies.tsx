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

interface Funnel {
  ID: number;
  name: string;
}

interface Company {
  ID: number;
  name: string;
  address: string;
  funnel_id?: number;
  funnel?: Funnel;
}

const companySchema = z.object({
  name: z.string().min(1, "Name is required"),
  address: z.string().optional(),
  funnel_id: z.number().optional(),
});

type CompanyFormValues = z.infer<typeof companySchema>;

export default function Companies() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [funnels, setFunnels] = useState<Funnel[]>([]);
  const [isCreating, setIsCreating] = useState(false);
  const [editingCompany, setEditingCompany] = useState<Company | null>(null);

  const form = useForm<CompanyFormValues>({
    resolver: zodResolver(companySchema),
    defaultValues: {
      name: '',
      address: '',
      funnel_id: 0,
    },
  });

  const fetchData = async () => {
      try {
        const [compRes, funnelRes] = await Promise.all([
          api.get('/api/companies'),
          api.get('/api/funnels')
        ]);
        setCompanies(compRes.data);
        setFunnels(funnelRes.data);
      } catch (error) {
        console.error("Failed to fetch data", error);
      }
    };

  useEffect(() => {
    fetchData();
  }, []);

  const onEdit = (company: Company) => {
      setEditingCompany(company);
      setIsCreating(true);
      form.reset({
          name: company.name,
          address: company.address,
          funnel_id: company.funnel_id || 0,
      });
  };

  const onSubmit = async (data: CompanyFormValues) => {
      try {
          const payload = { ...data, funnel_id: data.funnel_id === 0 ? null : data.funnel_id };

          if (editingCompany) {
              await api.put(`/api/companies/${editingCompany.ID}`, payload);
          } else {
              await api.post('/api/companies', payload);
          }
          await fetchData();
          setIsCreating(false);
          setEditingCompany(null);
          form.reset({ name: '', address: '', funnel_id: 0 });
      } catch (error) {
          console.error("Failed to save company", error);
      }
  };

  const handleCancel = () => {
      setIsCreating(false);
      setEditingCompany(null);
      form.reset({ name: '', address: '', funnel_id: 0 });
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold tracking-tight">Companies</h2>
        <Button onClick={() => { setIsCreating(true); setEditingCompany(null); form.reset({ name: '', address: '', funnel_id: 0 }); }}>
            {isCreating ? 'Cancel' : 'Add Company'}
        </Button>
      </div>
      
      {isCreating && (
          <Card>
              <CardHeader><CardTitle>{editingCompany ? 'Edit Company' : 'New Company'}</CardTitle></CardHeader>
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
                          
                          <FormField control={form.control} name="funnel_id" render={({ field }) => (
                              <FormItem>
                                  <FormLabel>Funnel</FormLabel>
                                  <FormControl>
                                      <select 
                                        className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                        {...field}
                                        value={field.value}
                                        onChange={(e) => field.onChange(parseInt(e.target.value, 10))}
                                      >
                                          <option value="0">Select a funnel (Optional)</option>
                                          {funnels.map(f => (
                                              <option key={f.ID} value={f.ID}>{f.name}</option>
                                          ))}
                                      </select>
                                  </FormControl>
                                  <FormMessage />
                              </FormItem>
                          )} />

                          <div className="flex space-x-2">
                              <Button type="submit">Save</Button>
                              <Button type="button" variant="outline" onClick={handleCancel}>Cancel</Button>
                          </div>
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
                        <TableHead>Funnel</TableHead>
                        <TableHead>Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {companies.map((company) => (
                        <TableRow key={company.ID}>
                            <TableCell>{company.ID}</TableCell>
                            <TableCell>{company.name}</TableCell>
                            <TableCell>{company.address}</TableCell>
                            <TableCell>{company.funnel?.name || '-'}</TableCell>
                            <TableCell>
                                <Button variant="ghost" size="sm" onClick={() => onEdit(company)}>Edit</Button>
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