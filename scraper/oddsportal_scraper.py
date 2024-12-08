import time
import csv
import os
import re
import io
import pandas as pd
import nltk
from selenium import webdriver
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from bs4 import BeautifulSoup
from datetime import datetime
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize
from transformers import BertTokenizer, BertForSequenceClassification,AutoTokenizer, AutoModelForSequenceClassification, pipeline
from zemberek import TurkishMorphology
from datetime import datetime, timezone
from ensemble_model import EnsembleModel
from flask import Flask,jsonify,request
import torch
import random
from models import RequestTweet,Tweet
import os

app = Flask(__name__)

# Gerekli dil işleme kütüphanelerini indir
nltk.download('punkt')
nltk.download('stopwords')
print("hello")
# Model ve tokenizer yolu
model1_path = r"/scraper/Models/final_model_pytorch"
model2_path = r"/scraper/Models/saved_model"
print(model1_path)
# Hugging Face pipeline modeli (Türkçe BERT modelini kullanıyoruz)
pipe = pipeline("sentiment-analysis", model="savasy/bert-base-turkish-sentiment-cased")

model1 = AutoModelForSequenceClassification.from_pretrained(model1_path)
model2 = AutoModelForSequenceClassification.from_pretrained(model2_path)

# Türkçe stopwords listesi
stop_words = set(stopwords.words('turkish'))

# BERT Tokenizer'ı yükle (bu kısım tokenization için eklenmiştir ancak şu an kullanılmıyor)
tokenizer = BertTokenizer.from_pretrained('dbmdz/bert-base-turkish-cased')

# Zemberek morfoloji analizörü
morphology = TurkishMorphology.create_with_defaults()

ensemble_model = EnsembleModel(model1, model2, pipe)


def clean_process(text):
    # 1. Özel karakterlerin, sayılardaki karakterlerin ve sembollerin kaldırılması
    text = re.sub(r'[^A-Za-zçÇğĞıİöÖşŞüÜ0-9\s]', '', text)
    
    # 2. Küçük harfe dönüştürme
    text = text.lower()
    
    # 3. Stopwords kaldırma
    filtered_words = [word for word in word_tokenize(text) if word not in stop_words]

    # 4. Zemberek ile kök bulma (stemming)
    # For simplicity, I'm skipping the stemming part in this version
    return ' '.join(filtered_words)


def clean_text(tweetSavePath,cleanedSavePath):
    file_path = tweetSavePath
    
    # tweetler.csv dosyasının varlığını kontrol et
    if not os.path.exists(file_path):
        print(f"{file_path} dosyası bulunamadı!")
        return []

    try:
        # tweetler.csv dosyasını pandas DataFrame olarak oku
        df = pd.read_csv(file_path, encoding='utf-8-sig')

        # "Tweet" ve "Date" sütunlarının olup olmadığını kontrol et
        if 'Tweet' not in df.columns or 'Date' not in df.columns:
            print('"Tweet" ve "Date" sütunları bulunamadı!')
            return []

        # Tweet sütununu temizle
        df['Cleaned_Tweet'] = df['Tweet'].apply(clean_process)
        
        # Date sütununu olduğu gibi koruyarak, sadece Cleaned_Tweet ve Date sütunlarını kaydet
        df[['Date', 'Cleaned_Tweet']].to_csv(cleanedSavePath, index=False, encoding='utf-8-sig')

        print(f"Temizlenmiş tweetler şu dosyaya kaydedildi: {cleanedSavePath}")
        return df[['Date', 'Tweet']]

    except Exception as e:
        print(f"Bir hata oluştu: {e}")
        return []

def launch_webdriver():
    chrome_options = webdriver.ChromeOptions()
    chrome_options.add_argument("--headless")  # Başsız mod
    chrome_options.add_argument("--disable-gpu")  # GPU kullanımı devre dışı
    chrome_options.add_argument("--window-size=1920,1080")  # Tam pencere boyutu
    chrome_options.add_argument("--disable-extensions")  # Eklentileri devre dışı bırak
    chrome_options.add_argument("--disable-dev-shm-usage")  # /dev/shm kullanımını azalt
    chrome_options.add_argument("--no-sandbox")  # Güvenli mod devre dışı
    chrome_options.add_argument("--ignore-certificate-errors")  # SSL hatalarını yoksay       
    return webdriver.Remote("http://selenium:4444/wd/hub", options=chrome_options)


def ProcessRequest(request_Tweet,tweetSavePath,cleanedSavePath,ClassifySavePath,driver):
    
    web_url = 'https://x.com/i/flow/login'
    driver.get(web_url)
    driver.implicitly_wait(50)
    

    input_xpath = '//div[@class="css-175oi2r r-18u37iz r-16y2uox r-1wbh5a2 r-1wzrnnt r-1udh08x r-xd6kpl r-is05cd r-ttdzmv"]//input[@name="text"]'
    input_element = driver.find_element(By.XPATH, input_xpath)
    input_element.send_keys("loncito123")  

    time.sleep(3)

    button_xpath1 = '//button[.//span[text()="Next"]]'
    login_button1 = driver.find_element(By.XPATH, button_xpath1)
    driver.execute_script("arguments[0].scrollIntoView();", login_button1)
    driver.execute_script("arguments[0].click();", login_button1)
    
    time.sleep(3)
   
    password_input_xpath = '//input[@name="password"]'
    password_input = driver.find_element(By.XPATH, password_input_xpath)
    password_input.send_keys("1u3JWdfhNS")  

    time.sleep(3)

    login_button_xpath = '//span[text()="Log in"]'
    login_button = driver.find_element(By.XPATH, login_button_xpath)
    driver.execute_script("arguments[0].scrollIntoView();", login_button)
    driver.execute_script("arguments[0].click();", login_button)
    
    time.sleep(3)

    element = driver.find_element(By.XPATH, '/html/body/div[1]/div/div/div[2]/header/div/div/div/div[1]/div[2]/nav/a[2]/div')
    element.click()
    
    time.sleep(3)

    search_box = WebDriverWait(driver, 60).until(
        EC.presence_of_element_located((By.XPATH, '//input[@data-testid="SearchBox_Search_Input"]'))
    )
    
    search_query= ""

    if isinstance(request_Tweet,RequestTweet):
        startedDateFormat =  datetime.strptime(request_Tweet.startDate, "%d/%m/%Y").strftime("%Y-%m-%d")
        endDateFormat =  datetime.strptime(request_Tweet.endDate, "%d/%m/%Y").strftime("%Y-%m-%d")
        search_query = f"{request_Tweet.tag} lang:tr until:{endDateFormat} since:{startedDateFormat}"



    search_box.send_keys(search_query)

    search_box.submit()   

    
    all_tweets = []
    seen_tweets = set()
    scroll_count = 0
    max_scrolls = 10  # Maksimum kaydırma sayısı

    while scroll_count < max_scrolls:
        tweet_elements = driver.find_elements(By.XPATH, '//article[@data-testid="tweet"]')

        if not tweet_elements:
                print("Yeni tweet bulunamadı, işlem tamamlandı.")
                break


        for tweet in tweet_elements:
            try:
                tweet_text = tweet.find_element(By.XPATH, './/div[@data-testid="tweetText"]').text
                tweet_date = tweet.find_element(By.XPATH, './/time').get_attribute('datetime')
                driver.save_screenshot("screenshot6.png")
                tweet_date = datetime.strptime(tweet_date, '%Y-%m-%dT%H:%M:%S.%fZ').replace(tzinfo=timezone.utc)  # datetime nesnesi oluşturuldu
                tweet_date = tweet_date.strftime('%d/%m/%Y')  #

                if tweet_text not in seen_tweets:
                    all_tweets.append((tweet_text, tweet_date))
                    seen_tweets.add(tweet_text)

            except Exception as e:
                print(f"Tweet işlenirken hata oluştu: {e}")

        # Sayfayı kaydır
        driver.execute_script("window.scrollBy(0, 800);")
        time.sleep(3)
        scroll_count += 1    

    try:
        
        output = io.StringIO()
        writer = csv.writer(output)
        writer.writerow(['Date', 'Tweet'])
        for tweet, date in all_tweets:
            writer.writerow([date, tweet])

        output.seek(0)  # Dosya başına geri dön
        # Geçici CSV dosyasını oluşturma
        with open(tweetSavePath, mode='w', newline='', encoding='utf-8-sig') as file:
             file.write(output.getvalue())  # CSV verisini dosyaya yaz
    except Exception as e:
        print(f"Hata oluştu: {e}")

    
    clean_text(tweetSavePath,cleanedSavePath)

def classify_csv(cleanedSavePath,ClassifySavePath,driver):

    data = pd.read_csv(cleanedSavePath, encoding='utf-8-sig')  
    predicted_labels = []

    for text in data["Cleaned_Tweet"]:
        final_category, final_probs = ensemble_model.predict(text)
        predicted_labels.append(ensemble_model.labels[final_category])  # Etiket adı ekleniyor

    # Tahmin edilen etiketleri veri çerçevesine ekle
    data["label"] = predicted_labels

    # Sonuçları CSV dosyasına kaydet
    data.to_csv(ClassifySavePath, index=False, encoding='utf-8-sig')
    print(f"Kategorize edilmiş veriler başarıysSla kaydedildi: {ClassifySavePath}")

    time.sleep(10)
    driver.quit()
def ConvertData(ClassifySavePath):
    data = pd.read_csv(ClassifySavePath, encoding='utf-8-sig')  
    tweet_list = [Tweet(row['Date'], row['Cleaned_Tweet'], row['label']) for _, row in data.iterrows()]
    result_list = [tweet.to_dict() for tweet in tweet_list]

    return result_list


@app.route('/getValue',methods=['POST'])
def GetResults():
    data = request.get_json()
    start = data.get('startedDate')
    end = data.get('endDate')
    hashTag = data.get('hashTag')
    category = data.get('category')
    limit = data.get('quantityLimit')
    print(category,limit)

    if not all([start, end, hashTag]):
        return jsonify({"error": "start, end ve tag parametreleri gerekli!"}), 400
    
    tweet_request  = RequestTweet(start,end,hashTag)
    prefixRand = int(''.join(random.choices('0123456789',k=6)))
    tweetSavePath = "/scraper/results/tweet_" + str(prefixRand) + ".csv"
    cleanedSavePath = "/scraper/results/cleaned_" + str(prefixRand) + ".csv"
    ClassifySavePath = "/scraper/results/classified_" + str(prefixRand) + ".csv"

    driver = launch_webdriver()

    print("Started Data",tweet_request.startDate,"End Date",tweet_request.endDate+"Tag:",tweet_request.tag)
    ProcessRequest(tweet_request ,tweetSavePath,cleanedSavePath,ClassifySavePath,driver)
    classify_csv(cleanedSavePath,ClassifySavePath,driver)
    result = ConvertData(ClassifySavePath)
    return jsonify({"tweets": result})


if __name__ == "__main__":
    app.run(debug=True,host='0.0.0.0',port=5000)
