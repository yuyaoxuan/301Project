import os
import csv
import random
from datetime import datetime, timedelta

# Configuration
# Function to generate clients starting from a given number
def generate_clients(start_number, num_clients=5):
    return [f'client{i}' for i in range(start_number, start_number + num_clients)]

# Start from client number 1 and generate 5 clients
start_number = 1
clients = generate_clients(start_number)
start_year = 2024
end_year = 2025
start_month = 12
end_month = 3
num_transactions_per_month = 50

# Generate the transactions
# Global counter for transaction IDs
global_last_id = 1000  # Starting point for global ID

# Function to generate a random transaction
def generate_transaction(last_id, client_id, current_date, time_increment):
    # Increment the ID
    new_id = last_id + 1
    
    # Vary the transaction type
    transaction_type = random.choice(['D', 'W'])
    
    # Vary the amount (positive, 2 decimal places)
    amount = round(random.uniform(10.00, 10000.00), 2)
    
    # Increment the date by a fixed time increment plus a random variation
    rand_variation = timedelta(minutes=random.randint(-30, 30))  # Random variation of ±30 minutes
    new_date = current_date + time_increment + rand_variation
    
    # Ensure the new date doesn't go before the current date (due to negative random variation)
    if new_date < current_date:
        new_date = current_date + timedelta(minutes=5)  # Minimum gap of 5 minutes
    
    # Vary the status with 70% Completed, 15% Pending, 15% Failed
    status_weights = ['Completed'] * 70 + ['Pending'] * 15 + ['Failed'] * 15
    status = random.choice(status_weights)
    
    # Return the new row as a dictionary
    return {
        'ID': new_id,
        'ClientID': client_id,
        'Transaction': transaction_type,
        'Amount': f'{amount:.2f}',
        'Date': new_date.strftime('%Y-%m-%dT%H:%M:%SZ'),  # ISO 8601 format
        'Status': status
    }


current_year = start_year
current_month = start_month

while (current_year, current_month) <= (end_year, end_month):
    for client_id in clients:
        # Use global_last_id instead of resetting for each client
        start_date = datetime(current_year, current_month, 1, 0, 0, 0)
        end_date = (datetime(current_year, current_month + 1, 1, 0, 0, 0) if current_month < 12 
                    else datetime(current_year + 1, 1, 1, 0, 0, 0))
        last_day_of_month = (end_date - timedelta(days=1)).replace(hour=23, minute=59, second=59)

        total_minutes = (end_date - start_date).total_seconds() / 60
        base_time_increment = timedelta(minutes=total_minutes // num_transactions_per_month)

        transactions = []
        current_date = start_date
        for _ in range(num_transactions_per_month):
            transaction = generate_transaction(global_last_id, client_id, current_date, base_time_increment)
            transactions.append(transaction)
            global_last_id = transaction['ID']  # Update global_last_id after each transaction
            current_date = datetime.strptime(transaction['Date'], '%Y-%m-%dT%H:%M:%SZ')

        # Ensure last transaction does not go past the last day of the month
        last_transaction_date = datetime.strptime(transactions[-1]['Date'], '%Y-%m-%dT%H:%M:%SZ')
        if last_transaction_date > last_day_of_month:
            transactions[-1]['Date'] = last_day_of_month.strftime('%Y-%m-%dT%H:%M:%SZ')

        # Create directory for client transactions
        output_directory = f'./transaction-logs/{client_id}/'
        os.makedirs(output_directory, exist_ok=True)

        # Define output filename
        output_filename = f'txn_log_{current_year}_{str(current_month).zfill(2)}.csv'
        output_file_path = os.path.join(output_directory, output_filename)

        # Save transactions to CSV
        with open(output_file_path, mode='w', newline='') as file:
            fieldnames = ['ID', 'ClientID', 'Transaction', 'Amount', 'Date', 'Status']
            writer = csv.DictWriter(file, fieldnames=fieldnames)
            writer.writeheader()
            writer.writerows(transactions)

        print(f"{num_transactions_per_month} transactions saved for {client_id} in {current_year}-{str(current_month).zfill(2)} at {output_file_path}")

    # Move to the next month
    if current_month == 12:
        current_year += 1
        current_month = 1
    else:
        current_month += 1
