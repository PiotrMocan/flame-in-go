import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import api from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';

interface Funnel {
  ID: number;
  name: string;
  next_funnels: Funnel[];
  previous_funnels: Funnel[];
}

const funnelSchema = z.object({
  name: z.string().min(1, "Name is required"),
  next_funnel_ids: z.array(z.number()).optional(),
  previous_funnel_ids: z.array(z.number()).optional(),
});

type FunnelFormValues = z.infer<typeof funnelSchema>;

export default function Funnels() {
  const [funnels, setFunnels] = useState<Funnel[]>([]);
  const [isCreating, setIsCreating] = useState(false);
  const [editingFunnel, setEditingFunnel] = useState<Funnel | null>(null);

  const form = useForm<FunnelFormValues>({
    resolver: zodResolver(funnelSchema),
    defaultValues: {
      name: '',
      next_funnel_ids: [],
      previous_funnel_ids: [],
    },
  });

  const fetchData = async () => {
    try {
      const funnelsRes = await api.get('/api/funnels');
      setFunnels(funnelsRes.data);
    } catch (error) {
      console.error("Failed to fetch data", error);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const onEdit = (funnel: Funnel) => {
      setEditingFunnel(funnel);
      setIsCreating(true);
      form.reset({
          name: funnel.name,
          next_funnel_ids: funnel.next_funnels ? funnel.next_funnels.map(f => f.ID) : [],
          previous_funnel_ids: funnel.previous_funnels ? funnel.previous_funnels.map(f => f.ID) : [],
      });
  };

  const onDelete = async (id: number) => {
      if (!confirm("Are you sure you want to delete this funnel?")) return;
      try {
          await api.delete(`/api/funnels/${id}`);
          fetchData();
      } catch (error) {
          console.error("Failed to delete funnel", error);
      }
  };

  const onSubmit = async (data: FunnelFormValues) => {
      try {
          if (editingFunnel) {
              await api.put(`/api/funnels/${editingFunnel.ID}`, data);
          } else {
              await api.post('/api/funnels', data);
          }
          await fetchData();
          setIsCreating(false);
          setEditingFunnel(null);
          form.reset({ name: '', next_funnel_ids: [], previous_funnel_ids: [] });
      } catch (error) {
          console.error("Failed to save funnel", error);
      }
  };

  const handleCancel = () => {
      setIsCreating(false);
      setEditingFunnel(null);
      form.reset({ name: '', next_funnel_ids: [], previous_funnel_ids: [] });
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold tracking-tight">Funnels</h2>
        <Button onClick={() => { setIsCreating(true); setEditingFunnel(null); form.reset({ name: '', next_funnel_ids: [], previous_funnel_ids: [] }); }}>
            Add Funnel
        </Button>
      </div>
      
      {isCreating && (
          <Card>
              <CardHeader><CardTitle>{editingFunnel ? 'Edit Funnel' : 'New Funnel'}</CardTitle></CardHeader>
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
                          
                          <div className="grid grid-cols-2 gap-4">
                            <FormField control={form.control} name="previous_funnel_ids" render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Previous Funnels</FormLabel>
                                    <div className="border p-4 rounded-md max-h-60 overflow-y-auto">
                                        {funnels.filter(f => f.ID !== editingFunnel?.ID).map(f => (
                                            <div key={f.ID} className="flex items-center space-x-2 mb-2">
                                                <input
                                                    type="checkbox"
                                                    id={`prev-funnel-${f.ID}`}
                                                    checked={field.value?.includes(f.ID)}
                                                    onChange={(e) => {
                                                        const checked = e.target.checked;
                                                        const current = field.value || [];
                                                        if (checked) {
                                                            field.onChange([...current, f.ID]);
                                                        } else {
                                                            field.onChange(current.filter(id => id !== f.ID));
                                                        }
                                                    }}
                                                    className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                                                />
                                                <Label htmlFor={`prev-funnel-${f.ID}`}>{f.name}</Label>
                                            </div>
                                        ))}
                                        {funnels.length <= (editingFunnel ? 1 : 0) && <span className="text-sm text-muted-foreground">No other funnels available.</span>}
                                    </div>
                                    <FormMessage />
                                </FormItem>
                            )} />

                            <FormField control={form.control} name="next_funnel_ids" render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Next Funnels</FormLabel>
                                    <div className="border p-4 rounded-md max-h-60 overflow-y-auto">
                                        {funnels.filter(f => f.ID !== editingFunnel?.ID).map(f => (
                                            <div key={f.ID} className="flex items-center space-x-2 mb-2">
                                                <input
                                                    type="checkbox"
                                                    id={`next-funnel-${f.ID}`}
                                                    checked={field.value?.includes(f.ID)}
                                                    onChange={(e) => {
                                                        const checked = e.target.checked;
                                                        const current = field.value || [];
                                                        if (checked) {
                                                            field.onChange([...current, f.ID]);
                                                        } else {
                                                            field.onChange(current.filter(id => id !== f.ID));
                                                        }
                                                    }}
                                                    className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                                                />
                                                <Label htmlFor={`next-funnel-${f.ID}`}>{f.name}</Label>
                                            </div>
                                        ))}
                                        {funnels.length <= (editingFunnel ? 1 : 0) && <span className="text-sm text-muted-foreground">No other funnels available.</span>}
                                    </div>
                                    <FormMessage />
                                </FormItem>
                            )} />
                          </div>

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
          <CardTitle>All Funnels</CardTitle>
        </CardHeader>
        <CardContent>
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>ID</TableHead>
                        <TableHead>Name</TableHead>
                        <TableHead>Previous Funnels</TableHead>
                        <TableHead>Next Funnels</TableHead>
                        <TableHead>Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {funnels.map((funnel) => (
                        <TableRow key={funnel.ID}>
                            <TableCell>{funnel.ID}</TableCell>
                            <TableCell>{funnel.name}</TableCell>
                            <TableCell>
                                {funnel.previous_funnels?.map(pf => pf.name).join(', ') || '-'}
                            </TableCell>
                            <TableCell>
                                {funnel.next_funnels?.map(nf => nf.name).join(', ') || '-'}
                            </TableCell>
                            <TableCell className="space-x-2">
                                <Button variant="ghost" size="sm" onClick={() => onEdit(funnel)}>Edit</Button>
                                <Button variant="destructive" size="sm" onClick={() => onDelete(funnel.ID)}>Delete</Button>
                            </TableCell>
                        </TableRow>
                    ))}
                    {funnels.length === 0 && (
                        <TableRow>
                            <TableCell colSpan={5} className="text-center h-24 text-muted-foreground">No funnels found.</TableCell>
                        </TableRow>
                    )}
                </TableBody>
            </Table>
        </CardContent>
      </Card>
    </div>
  );
}
