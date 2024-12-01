from flask import Flask, request, jsonify

app = Flask(__name__)

class AIRequest:
    def __init__(self, started_date, end_date, hashtag, category, quantity_limit):
        self.started_date = started_date
        self.end_date = end_date
        self.hashtag = hashtag
        self.category = category
        self.quantity_limit = quantity_limit

    def __repr__(self):
        return f"AIRequest(started_date={self.started_date}, end_date={self.end_date}, hashtag={self.hashtag}, category={self.category}, quantity_limit={self.quantity_limit})"


@app.route('/endpoint', methods=['POST'])
def ai_request():

    data = request.get_json()
    
    started_date = data.get('startedDate')
    end_date = data.get('endDate')
    hashtag = data.get('hashTag')
    category = data.get('category')
    quantity_limit = data.get('quantityLimit')
    
    ai_request = AIRequest(
        started_date=started_date,
        end_date=end_date,
        hashtag=hashtag,
        category=category,
        quantity_limit=quantity_limit
    )
        
    # Yanıt olarak 'OK' dön
    return jsonify({"message": "OK"}), 200

if __name__ == '__main__':
    app.run(debug=True)
